package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/ahmetfurkanunal/bitaksi-taxihub/db"
	"github.com/ahmetfurkanunal/bitaksi-taxihub/driver"

	docs "github.com/ahmetfurkanunal/bitaksi-taxihub/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	mongoURI := "mongodb://mongo:27017"

	if mongoURI == "" {
		log.Fatal(" MONGO_URI is NOT defined!")
	}

	client, err := db.ConnectMongo(mongoURI)
	if err != nil {
		log.Fatalf("mongo connection error: %v", err)
	}

	defer func(c *mongo.Client) {
		if err := c.Disconnect(nil); err != nil {
			log.Printf("mongo disconnect error: %v", err)
		}
	}(client)

	log.Println(" MongoDB connected to:", mongoURI)
	log.Println(" Active DB: taxihub")

	router := gin.Default()

	corsCfg := cors.Config{
		AllowAllOrigins: true,
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
			"X-Requested-With",
		},
	}
	router.Use(cors.New(corsCfg))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	driverRepo := driver.NewRepository(client)
	driverService := driver.NewService(driverRepo)
	driverHandler := driver.NewHandler(driverService)
	driverHandler.RegisterRoutes(router)

	docs.SwaggerInfo.BasePath = "/"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println("API started on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
