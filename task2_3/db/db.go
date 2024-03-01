package db

import (
	"context"
	"fmt"
	"log"
)

type Result struct {
	Category
	Products []Product
}

func (pg *Postgres) GetProductById(id int) (Product, error) {
	var (
		product       Product
		categories_id int64
		category      Category
	)

	tx, err := pg.db.Begin(context.Background())
	if err != nil {
		log.Println("Error to connect db")
		return Product{}, err
	}

	defer tx.Rollback(context.Background())

	err = tx.QueryRow(
		context.Background(),
		`SELECT name, mark FROM testschema.products WHERE id = $1`, id).Scan(
		&product.Name, &product.Mark)

	if err != nil {
		log.Printf("Failure to select product %+v", err)
		return product, err
	}

	allCategories, err := tx.Query(context.Background(),
		`SELECT categories_id FROM testschema.products_categories WHERE products_id = $1`, id)

	if err != nil {
		log.Printf("Error to select categories: %+v", err)
	}

	defer allCategories.Close()

	conn, err := pg.db.Begin(context.Background())

	if err != nil {
		log.Printf("Error to create connection pool: %+v", err)
	}

	defer conn.Rollback(context.Background())

	for allCategories.Next() {
		if err = allCategories.Scan(&categories_id); err != nil {
			log.Println("Error to search categories")
			return Product{}, err
		}

		err = conn.QueryRow(context.Background(),
			`SELECT name FROM testschema.categories WHERE id = $1`,
			categories_id).Scan(&category.Name)

		if err != nil {
			log.Printf("Error to select category: %+v", err)
			return Product{}, err
		}

		product.Categories = append(product.Categories, category)
	}

	return product, nil
}

func (pg *Postgres) GetProductInCategory(arg string) ([]Result, error) {
	var (
		categories_id   int64
		categories_name string
		resultArray     []Result
		product         Product
	)

	tx, err := pg.db.Begin(context.Background())
	if err != nil {
		log.Println("Error to connect db")
		return nil, err
	}

	defer tx.Rollback(context.Background())

	allCategories, err := tx.Query(context.Background(),
		`SELECT * FROM testschema.categories`)

	if err != nil {
		log.Printf("Error to select categories: %+v", err)
	}

	defer allCategories.Close()

	conn, err := pg.db.Begin(context.Background())

	if err != nil {
		log.Printf("Error to create connection pool: %+v", err)
	}

	defer conn.Rollback(context.Background())

	for allCategories.Next() {
		productsArray := make([]Product, 0)
		if err = allCategories.Scan(&categories_id, &categories_name); err != nil {
			log.Println("Error to search categories")
			return nil, err
		}

		query := CheckCondition(arg, categories_id)

		products, err := conn.Query(context.Background(), query)

		if err != nil {
			log.Printf("Error to select products: %+v", err)
			return nil, err
		}

		for products.Next() {
			if err = products.Scan(&product.Name, &product.Mark); err != nil {
				log.Printf("Error to parse products: %+v", err)
				return nil, err
			}

			productsArray = append(productsArray, product)
		}

		resultArray = append(resultArray, Result{
			Category: Category{Name: categories_name},
			Products: productsArray,
		})
	}

	return resultArray, nil
}

func CheckCondition(condition string, id int64) string {
	var query string
	if condition == "MAX" {
		query = fmt.Sprintf(`
					SELECT prod.name, prod.mark
					FROM testschema.products AS prod
					JOIN testschema.products_categories AS pc ON 
					prod.id = pc.products_id
					WHERE pc.categories_id = %d AND prod.mark = (
						SELECT MAX(prod.mark) FROM testschema.products AS prod
						JOIN testschema.products_categories AS pc ON
						prod.id = pc.products_id
						WHERE pc.categories_id = %d
					)`, id, id)
	} else {
		query = fmt.Sprintf(`
					SELECT prod.name, prod.mark
					FROM testschema.products AS prod
					JOIN testschema.products_categories AS pc ON 
					prod.id = pc.products_id
					WHERE pc.categories_id = %d AND prod.mark = (
						SELECT MIN(prod.mark) FROM testschema.products AS prod
						JOIN testschema.products_categories AS pc ON
						prod.id = pc.products_id
						WHERE pc.categories_id = %d
					)`, id, id)
	}

	return query
}

//SELECT cat.name, prod.name, prod.mark
//	FROM testschema.categories as cat
//	join testschema.products_categories as pc on
//	cat.id = pc.categories_id
//	join testschema.products as prod on
//	pc.products_id = prod.id
//	where prod.mark = (
//		select MIN(mark)
//		from testschema.products as prod2
//		join testschema.products_categories as pc2 on
//		prod2.id = pc2.products_id
//		where pc2.categories_id = pc.categories_id
//	)
//order by cat.name
