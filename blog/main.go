package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

const (
	DBHOST = "127.0.0.1"
	DBPORT = ":3306"
	DBUser = "exampleuser"
	DBPass = "examplepass"
	DBDbase = "exampledb"
	PORT =":8443"
)

var database *sql.DB


func (p Page) TruncatedText () string {
	chars := 0
	for i := range p.Content {
		chars++
		if chars > 100 {
			return p.RawContent[:i] + `...`
		}
	}
	return p.RawContent 
}

func ServePage (c *gin.Context) {
	pageGUID := c.Param("guid")
	thisPage := Page{}
	fmt.Println(pageGUID)
	err := database.QueryRow("SELECT page_title,page_content,page_date FROM pages WHERE page_guid=?", pageGUID).Scan(&thisPage.Title, &thisPage.RawContent, &thisPage.Date)
	thisPage.Content = template.HTML(thisPage.RawContent)
	if err != nil {
		log.Println("Couldn't get page: "+pageGUID)
		log.Println(err.Error())
		c.Redirect(http.StatusFound, "/not-found")
	}

	c.HTML(http.StatusOK, "blog.tmpl", gin.H{
		"Title": thisPage.Title,
		"Date": thisPage.Date,
		"Content": thisPage.Content,
	})
}

func NotFound (c *gin.Context) {
		c.HTML(http.StatusNotFound, "not-found.tmpl", gin.H{})
}

func RedirIndex(c *gin.Context){
		c.Redirect(http.StatusMovedPermanently, "/home")
}

func ServerIndex(c *gin.Context){
	var Pages = []Page{}
	query := "SELECT page_title,page_content,page_date,page_guid FROM pages ORDER BY page_date DESC"
	pages, err := database.Query(query)
	if err != nil {
		c.Error(err)
	}
	defer pages.Close()
	for pages.Next() {
		thisPage := Page{}
		err := pages.Scan(&thisPage.Title, &thisPage.RawContent, &thisPage.Date, &thisPage.GUID)
			if err != nil {
				log.Println("Could not get query page!")
				log.Println(err.Error())
			}
		thisPage.Content = template.HTML(thisPage.RawContent)
		Pages = append(Pages, thisPage)
	}
	c.HTML(http.StatusOK, "home.tmpl", gin.H{"Pages": Pages})
}

func main()  {
	dbConn := fmt.Sprintf("%s:%s@tcp(%s)/%s", DBUser, DBPass, DBHOST, DBDbase)
	db, err := sql.Open("mysql", dbConn)
	if err != nil {
		log.Println("Couldn't connect to the database!")
		log.Println(err.Error())
	}

	database = db

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/", RedirIndex)
	router.GET("/home", ServerIndex)

	router.GET("/page/:guid", ServePage)

	router.GET("/not-found", NotFound)
	router.NoRoute(NotFound)

	RouteAPI(router)

	// HTTP
	// router.Run(PORT)

	// HTTPS
	log.Fatal(router.RunTLS(PORT, "server.crt", "server.key"))
}