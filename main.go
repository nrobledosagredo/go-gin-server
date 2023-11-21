package main

import (
	"context"
	"log"
	"os"

	"geoar-backend/routes"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Log
	logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logFile)

	// Conexión a MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Success connection to MongoDB.")
	}

	// Configuración de Gin
	r := gin.Default()

	// Configurar rutas
	routes.PointsRoutes(r, client)
	// Agregar la configuración de otras rutas

	// Iniciar servidor
	r.Run(":8080")
}
