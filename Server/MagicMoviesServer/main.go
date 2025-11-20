package main

import (
	"fmt"

	controllers "github.com/Swarnimrajsanu/MagicMoviesStream/Server/MagicMoviesServer/controllers"
	"github.com/Swarnimrajsanu/MagicMoviesStream/Server/MagicMoviesServer/database"
	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	router.GET("/hello", func(c *gin.Context) {
		c.String(200, "Hello, Magic Movies!")
	})

	router.GET("/movies", controllers.GetMovies(database.DBInstance()))
	router.GET("/movie/:imdb_id", controllers.GetMovie(database.DBInstance()))
	router.POST("/addmovie", controllers.AddMovie(database.DBInstance()))
	router.POST("/register", controllers.RegisterUser())
	router.POST("/login", controllers.LoginUser(database.DBInstance()))

	if err := router.Run(":8080"); err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
