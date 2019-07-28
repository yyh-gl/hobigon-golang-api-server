package handler

import (
	"github.com/julienschmidt/httprouter"
)

var Router *httprouter.Router

func init() {
	router := httprouter.New()

	router.POST("/api/v1/tasks", NotifyTaskHandler)

	Router = router
}
