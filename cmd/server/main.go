package main

import (
	"fmt"
	"log"

	"github.com/DevNCH/NAS/internal/database"
)

func main() {

	if err := database.Connect(); err != nil {
    	log.Fatal(err)
	}

	if err := database.Migrate(); err != nil {
    	log.Fatal(err)
	}

	fmt.Println("Servidor iniciado")
}