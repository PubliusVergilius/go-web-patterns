package main

import "html/template"

type Page struct {
	Title string
	RawContent string
	Content template.HTML
	Date string
	GUID string
	Session Session
}

type JSONResponse struct {
	Fields map[string]string
}