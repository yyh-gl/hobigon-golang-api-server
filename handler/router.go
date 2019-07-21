package handler

import (
	"github.com/julienschmidt/httprouter"
)

var Router *httprouter.Router

func init() {
	router := httprouter.New()

	router.GET("/api/v1/books", BookIndex)
	router.POST("/api/v1/tasks", TaskHandler)

	Router = router
}
