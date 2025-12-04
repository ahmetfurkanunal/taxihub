// @title TaxiHub Driver Service API
// @version 1.0
// @description Sürücü kayıtlarını yöneten ve yakın sürücüleri bulan TaxiHub Driver Service.
// @host localhost:8080
// @BasePath /
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ahmetfurkanunal/bitaksi-taxihub/db"
	"github.com/ahmetfurkanunal/bitaksi-taxihub/driver"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	docs "github.com/ahmetfurkanunal/bitaksi-taxihub/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {

	mongoURI := os.Getenv("MONGO_URI")

	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	client, err := db.ConnectMongo(mongoURI)
	if err != nil {
		log.Fatalf("mongo connection error: %v", err)
	}
	defer func(client *mongo.Client) {
		if err := client.Disconnect(nil); err != nil {
			log.Printf("mongo disconnect error: %v", err)
		}
	}(client)

	log.Println("MongoDB connected")

	router := gin.Default()

	docs.SwaggerInfo.BasePath = "/"

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	driverRepo := driver.NewRepository(client)
	driverService := driver.NewService(driverRepo)
	driverHandler := driver.NewHandler(driverService)
	driverHandler.RegisterRoutes(router)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
