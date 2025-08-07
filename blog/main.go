package main

import (
	"database/sql"
	"fmt"
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

func NotFound (c *gin.Context) {
		c.HTML(http.StatusNotFound, "not-found.tmpl", gin.H{})
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

	

	router.NoRoute(NotFound)
	RoutePages(router)
	RouteAPI(router)

	// HTTP
	// router.Run(PORT)

	// HTTPS
	log.Fatal(router.RunTLS(PORT, "server.crt", "server.key"))
}