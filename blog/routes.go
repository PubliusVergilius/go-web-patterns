package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

/********************* Page Routes *******************/

func (p Page) TruncatedText() string {
	chars := 0
	for i := range p.Content {
		chars++
		if chars > 100 {
			return p.RawContent[:i] + `...`
		}
	}
	return p.RawContent
}

func getPageComments(pageID int) (*[]Comment, error) {
	query := "SELECT id, comment_name as Name, comment_email, comment_text FROM comments WHERE page_id=?"
	_comments, err := database.Query(query, pageID)
	if err != nil {
		log.Println("Couldn't get comments to posts")
		log.Println(err)
		return nil, err
	}

	var comments []Comment
	for _comments.Next() {
		var comment Comment
		_comments.Scan(&comment.ID, &comment.Name, &comment.Email, &comment.CommentText)
		comments = append(comments, comment)
	}

	return &comments, nil
}

func ServePage(c *gin.Context) {
	pageGUID := c.Param("guid")
	thisPage := Page{}
	err := database.QueryRow("SELECT id,page_title,page_content,page_date FROM pages WHERE page_guid=?", pageGUID).Scan(&thisPage.ID, &thisPage.Title, &thisPage.RawContent, &thisPage.Date)
	thisPage.Content = template.HTML(thisPage.RawContent)
	if err != nil {
		log.Println("Couldn't get page: " + pageGUID)
		log.Println(err.Error())
		c.Redirect(http.StatusFound, "/not-found")
	}

	comments, err := getPageComments(thisPage.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not get comments!",
			"error":   err.Error(),
		})
	}
	log.Printf("Comments: %q", *comments)
	thisPage.Comments = *comments

	c.HTML(http.StatusOK, "blog.tmpl", gin.H{
		"Title":    thisPage.Title,
		"Date":     thisPage.Date,
		"Content":  thisPage.Content,
		"GUID":     thisPage.GUID,
		"Comments": thisPage.Comments,
	})
}

func RedirIndex(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "/home")
}

func ServerIndex(c *gin.Context) {
	var Pages = []Page{}
	query := "SELECT page_title,page_content,page_date,page_guid FROM pages ORDER BY page_date DESC"
	pages, err := database.Query(query)
	if err != nil {
		c.Error(err)
	}
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

func RoutePages(router *gin.Engine) {
	{
		pages := router.Group("/")
		pages.GET("/", RedirIndex)
		pages.GET("/home", ServerIndex)

		pages.GET("/page/:guid", ServePage)

		pages.GET("/not-found", NotFound)
	}
}

/********************* API Routes *******************/

func APIPage(c *gin.Context) {
	pageGUID := c.Param("guid")
	thisPage := Page{}

	query := "SELECT page_title,page_content,page_date FROM pages WHERE page_guid=?"
	err := database.QueryRow(query, pageGUID).Scan(&thisPage.Title, &thisPage.RawContent, &thisPage.Date)
	thisPage.Content = template.HTML(thisPage.RawContent)
	if err != nil {
		c.Error(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"title":   thisPage.Title,
		"content": thisPage.Content,
		"date":    thisPage.Date,
	})

}

type CommentForm struct {
	Name     string `form:"name" binding:"required"`
	Email    string `form:"email" binding:"required"`
	Comments string `form:"comments" binding:"required"`
	GUID     string `form:"guid" binding:"required"`
}

func APICommentPost(c *gin.Context) {
	// var commendAdded bool
	var form CommentForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// query := "INSERT INTO comments SET comment_name=?, comment_email=?, comment_text=?"
	query := "INSERT INTO comments (comment_name, comment_email, comment_text, page_id) SELECT ?, ?, ?, id FROM pages WHERE page_guid = ?;"
	res, err := database.Exec(query, form.Name, form.Email, form.Comments, form.GUID)
	if err != nil {
		log.Panicln(err.Error())
	}

	id, err := res.LastInsertId()
	if err != nil {
		// commendAdded = false
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error on saving comment!"})
		return
	}

	log.Printf("form added: %q", form)
	/// commendAdded = true

	c.JSON(http.StatusOK, gin.H{"id": id, "message": "Comment successfully posted!"})
}

func RouteAPI(router *gin.Engine) {
	{
		api := router.Group("/api")
		api.GET("/doc", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Olá!"})
		})
		api.GET("/pages/:guid", APIPage)

		api.POST("/comments", APICommentPost)
	}
}
