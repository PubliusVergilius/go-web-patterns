package main

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

func APIPage (c *gin.Context) {
	pageGUID := c.Param("guid")	
	thisPage := Page{}

	query := "SELECT page_title,page_content,page_date FROM pages WHERE page_guid=?"
	err := database.QueryRow(query, pageGUID).Scan(&thisPage.Title, &thisPage.RawContent, &thisPage.Date)
	thisPage.Content = template.HTML(thisPage.RawContent)
	if err != nil {
		c.Error(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"title": thisPage.Title,
		"content": thisPage.Content,
		"date": thisPage.Date,
	})

}

func RouteAPI (router *gin.Engine) {
	{
		api := router.Group("/api")
		api.GET("/doc", func (c *gin.Context){
			c.JSON(http.StatusOK, gin.H{"message": "Olá!"})
		})
		api.GET("/pages/:guid", APIPage)
	}	
}