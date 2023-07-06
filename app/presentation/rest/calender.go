package rest

import (
	"fmt"
	"github.com/yyh-gl/hobigon-golang-api-server/app"
	"net/http"
)

type Calender interface {
	Create(w http.ResponseWriter, r *http.Request)
}

type calender struct{}

func NewCalender() Calender {
	return &calender{}
}

// Create : カレンダー情報を新規作成
func (c calender) Create(w http.ResponseWriter, r *http.Request) {
	type (
		request struct {
			File string `validate:"required"`
		}
		response struct {
			Title string `json:"title"`
			Count int    `json:"count"`
		}
	)

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		app.Error(fmt.Errorf("error in http.Request.FormFile(): %w", err))
		DoResponse(w, err, http.StatusInternalServerError)
	}

	fmt.Println("========================")
	fmt.Println(file)
	fmt.Println("========================")
	fmt.Println(fileHeader)
	fmt.Println("========================")

	DoResponse(w, nil, http.StatusCreated)
}
