package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/teris-io/shortid"
)

var ctx = context.Background()

func main() {
	r := gin.Default()
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	r.GET("/:url", func(c *gin.Context) {
		resolvedLink, err := client.Get(ctx, c.Param("url")).Result()
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Can't resolve link with given id, link doesn't exist or expired.",
			})
		}
		c.Redirect(http.StatusFound, resolvedLink)
	})
	r.GET("/add-url", func(c *gin.Context) {
		sid, generatorCreatingError := shortid.New(1, shortid.DefaultABC, 2342)
		if generatorCreatingError != nil {
			panic(generatorCreatingError)
		}
		generatedSid, generatingError := sid.Generate()
		if generatingError != nil {
			panic(generatingError)
		}
		err := client.Set(ctx, generatedSid, c.Query("url"), 5*time.Minute).Err()
		if err != nil {
			panic(err)
		}
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Url shortened at %s", generatedSid),
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
