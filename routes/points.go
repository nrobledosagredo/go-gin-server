// Importar bibliotecas necesarias
package routes

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Función para extraer puntos de una colección de MongoDB según un filtro y una proyección.
func fetchPointsData(collection *mongo.Collection, filter bson.M, projection bson.M) ([]bson.M, error) {
	// Inicializar el slice para almacenar los resultados
	var results []bson.M
	// Establecer las opciones de la consulta para limitar los campos retornados
	opts := options.Find().SetProjection(projection)
	// Ejecutar la consulta en MongoDB
	cursor, err := collection.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, err
	}
	// Cerrar el cursor después de su uso
	defer cursor.Close(context.Background())
	// Iterar a través del cursor para recolectar todos los documentos
	for cursor.Next(context.Background()) {
		var singleResult bson.M
		if err := cursor.Decode(&singleResult); err != nil {
			return nil, err
		}
		results = append(results, singleResult)
	}
	// Retornar los resultados
	return results, nil
}

// Controlador para obtener todos los puntos geográficos
func getAllPoints(context *gin.Context, mongoClient *mongo.Client) {
	// Referencia a la colección "points" de la base de datos "GeoAR"
	pointsCollection := mongoClient.Database("GeoAR").Collection("points")
	// Consulta para obtener todos los puntos
	results, err := fetchPointsData(pointsCollection, bson.M{}, bson.M{})
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		log.Println("getAllPoints.")
	}
	// Retornar los resultados en formato JSON
	context.JSON(http.StatusOK, results)
}

// Controlador para obtener puntos geográficos por una ruta específica ("trail")
func getPointsByTrail(context *gin.Context, mongoClient *mongo.Client) {
	// Extraer el identificador de la ruta ("trail") desde el parámetro de la URL
	trail := context.Param("trail")
	// Convertir el identificador a un ObjectID válido de MongoDB
	objectID, err := primitive.ObjectIDFromHex(trail)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	// Referencia a la colección "points" de la base de datos "GeoAR"
	pointsCollection := mongoClient.Database("GeoAR").Collection("points")
	// Consulta para obtener puntos que coincidan con el identificador de la ruta
	results, err := fetchPointsData(pointsCollection, bson.M{"trail": objectID}, bson.M{"geometry": 1, "order": 1, "_id": 0})
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		log.Println("getPointsByTrails.")
	}
	// Retornar los resultados en formato JSON
	context.JSON(http.StatusOK, results)
}

// Configuración de las rutas HTTP para los puntos geográficos
func PointsRoutes(router *gin.Engine, mongoClient *mongo.Client) {
	// Log
	logFile, err := os.OpenFile("./app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logFile)

	// Ruta para obtener todos los puntos
	router.GET("/points", func(context *gin.Context) { getAllPoints(context, mongoClient) })
	// Ruta para obtener puntos por una ruta específica ("trail")
	router.GET("/points/:trail", func(context *gin.Context) { getPointsByTrail(context, mongoClient) })
}
