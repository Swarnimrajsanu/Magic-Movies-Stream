package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/Swarnimrajsanu/MagicMoviesStream/Server/MagicMoviesServer/database"
	"github.com/Swarnimrajsanu/MagicMoviesStream/Server/MagicMoviesServer/models"
	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var validate = validator.New()

func GetMovies(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		var movieCollection *mongo.Collection = database.OpenCollection("movies", client)

		cursor, err := movieCollection.Find(ctx, bson.D{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movies."})
		}
		defer cursor.Close(ctx)

		var movies []models.Movie

		if err = cursor.All(ctx, &movies); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode movies."})
			return
		}

		c.JSON(http.StatusOK, movies)

	}
}
func GetMovie(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {

		// FIXED: Use request context, not Gin context
		ctx, cancel := context.WithTimeout(c.Request.Context(), 100*time.Second)
		defer cancel()

		movieID := c.Param("imdb_id")
		if movieID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Movie ID is required"})
			return
		}

		movieCollection := database.OpenCollection("movies", client)

		var movie models.Movie

		err := movieCollection.FindOne(ctx, bson.D{{Key: "imdb_id", Value: movieID}}).Decode(&movie)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
			return
		}

		c.JSON(http.StatusOK, movie)
	}
}
func AddMovie(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {

		// FIX: Use request context
		ctx, cancel := context.WithTimeout(c.Request.Context(), 100*time.Second)
		defer cancel()

		var movie models.Movie

		if err := c.ShouldBindJSON(&movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		if err := validate.Struct(movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation failed",
				"details": err.Error(),
			})
			return
		}

		movieCollection := database.OpenCollection("movies", client)

		result, err := movieCollection.InsertOne(ctx, movie)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add movie"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Movie added successfully",
			"id":      result.InsertedID,
		})
	}
}
