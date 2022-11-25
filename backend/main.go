package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/keshav743/cinephile/database"
	"github.com/keshav743/cinephile/router"
)

func main() {
	loadEnvVariables()

	app := fiber.New()
	database.Connect()

	router.SetupRoutes(app)

	fmt.Println("Server up and running at port 3000.")
	
	log.Fatal(app.Listen(":3000"))
}

func loadEnvVariables() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
}