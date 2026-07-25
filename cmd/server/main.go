package main

import (
	"fmt"
	"log"

	"github.com/DevNCH/NAS/internal/database"
	"github.com/DevNCH/NAS/internal/routes"
)

func main() {

	if err := database.Connect(); err != nil {
    	log.Fatal(err)
	}

	if err := database.Migrate(); err != nil {
    	log.Fatal(err)
	}

	if err := database.Seed(); err != nil {
		log.Fatal(err)
	}

	router := routes.SetupRouter()

	fmt.Println("Servidor iniciado na porta 8080")

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Servidor iniciado")
}