package main

import (
	"context"
	"log"
	"main/config"
	"main/db"
)

func main() {
	cfg := config.ReadCfg()
	database, err := db.InitPG(context.Background(), *cfg)
	if err != nil {
		log.Fatalf("Error to create new pool, %+v", err)
		return
	}

}
