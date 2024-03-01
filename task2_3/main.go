package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"main/config"
	"main/db"
	"os"
	"strconv"
)

func printProduct(data db.Product) {
	fmt.Printf("Name: %s\nMark: %.1f\nCategories:\n", data.Name, data.Mark)
	for i, v := range data.Categories {
		fmt.Printf("\t%d. %s\n", i, v.Name)
	}
}

func printCategoryRating(data []db.Result) {
	for _, v := range data {
		fmt.Printf("Name category: %s\nProducts:\n", v.Category.Name)

		for _, val := range v.Products {
			fmt.Printf("\tName: %s\tMark: %0.1f\n", val.Name, val.Mark)
		}
	}
}

func menu(db *db.Postgres) {
	for {
		fmt.Printf("\nSelect option: \n" +
			"1.Find item by id\n2.Find max or min item\n3.Quit\n")
		scanner := bufio.NewScanner(os.Stdin)
		scanned := scanner.Scan()
		if !scanned {
			log.Fatal("program error")
			return
		}

		input := scanner.Text()

		switch input {
		case "1":
			fmt.Println("Input id")
			scanner.Scan()

			input = scanner.Text()
			id, err := strconv.Atoi(input)
			if err != nil {
				fmt.Println("error convert string to int")
				continue
			}

			data, err := db.GetProductById(id)

			if err != nil {
				fmt.Println(err)
				continue
			}

			printProduct(data)
			continue
		case "2":
			for {
				fmt.Printf("Select options:\n1.Max\n2.Min\n3.Return\n")

				scanner.Scan()
				input = scanner.Text()

				switch input {
				case "1":
					data, err := db.GetProductInCategory("MAX")
					if err != nil {
						fmt.Println(err)
						break
					}

					printCategoryRating(data)
					continue
				case "2":
					data, err := db.GetProductInCategory("MIN")
					if err != nil {
						fmt.Println(err)
						break
					}

					printCategoryRating(data)
					continue
				case "3":
				default:
					continue
				}
				break
			}
		case "3":
			fmt.Println("Goodbye")
			return
		default:
			continue
		}
	}
}

func main() {
	cfg := config.ReadCfg()
	database, err := db.InitPG(context.Background(), *cfg)
	if err != nil {
		log.Fatalf("Error to create new pool, %+v", err)
		return
	}

	menu(database)
}
