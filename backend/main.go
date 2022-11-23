package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	loadEnvVariables()
	MONGO_URI := "mongodb+srv://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@cluster0.iflqo.mongodb.net/?retryWrites=true&w=majority"

	app := fiber.New()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(MONGO_URI))
	handleDbError(err)

	defer client.Disconnect(ctx)

	fmt.Println("DB connection successfull")

	fmt.Println("Server up and running at port 3000.")

	app.Get("/", func (c *fiber.Ctx) error {
		err := c.SendString("And the API is UP!")
		fmt.Println(err)
        return err
	})
	
	log.Fatal(app.Listen(":3000"))

}

func loadEnvVariables() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
}

func handleDbError(err error) {
	if err != nil {
		log.Fatalf("DB connection failed. Err: %s", err)
	}
}
