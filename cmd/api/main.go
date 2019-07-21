package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/yyh-gl/hobigon-golang-api-server/handler"
)

func main() {
	fmt.Println("========================")
	fmt.Println("Server Start >> http://localhost:3000")
	fmt.Println("========================")
	log.Fatal(http.ListenAndServe(":3000", handler.Router))
}
