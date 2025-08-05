package main

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type File struct {
	ID string `uri:"id" binding:"required,validID"`
}

func validID(fl validator.FieldLevel) bool {
	id := fl.Field().String()
	fmt.Println(id)
	matched, _ := regexp.MatchString(`^[0-9]+$`, id)
	return matched
}

func main () {
	routes := gin.Default()

	routes.GET("/ping", func (c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	routes.GET("/pages", func (c *gin.Context) {
		filePath := "./var/www/static.html"
		c.File(filePath)
	})

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("validID", validID)
	}

	routes.GET("/pages/:id", func (c *gin.Context) {
		var file File
		if err := c.ShouldBindUri(&file); err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
			return 
		}
		filePath := "./var/www/"+file.ID+".html"
		c.File(filePath)
	})

	routes.Run()
}