package main

import "html/template"

type Page struct {
	ID         int
	Title      string
	RawContent string
	Content    template.HTML
	Comments   []Comment
	Date       string
	GUID       string
	// Session    Session
}

type JSONResponse struct {
	Fields map[string]string
}

type Comment struct {
	ID          int
	Name        string
	Email       string
	CommentText string
}
