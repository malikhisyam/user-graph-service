package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	relations "github.com/malikhisyam/user-graph-service/domains/relations/entities"
	users "github.com/malikhisyam/user-graph-service/domains/users/entities"
	"github.com/malikhisyam/user-graph-service/wizards"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	fmt.Println("Loading environment successfully....")
	wizards.PostgresDatabase.GetInstance().AutoMigrate(
		&users.User{},
		&relations.Follows{},
	)

	router := gin.Default()
	wizards.RegisterServer(router)
	router.Run(fmt.Sprintf(":%d", wizards.Config.Server.Port))
}