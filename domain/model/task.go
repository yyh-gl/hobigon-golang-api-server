package model

import "github.com/adlio/trello"

type Task struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Board struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Desc           string `json:"desc"`
	Closed         bool   `json:"closed"`
	IDOrganization string `json:"idOrganization"`
	Pinned         bool   `json:"pinned"`
	URL            string `json:"url"`
	ShortURL       string `json:"shortUrl"`
}

func ConvertToBoardModel(trelloBoard trello.Board) (board *Board) {
	board = new(Board)
	board.ID = trelloBoard.ID
	board.Name=  trelloBoard.Name
	board.Desc = trelloBoard.Desc
	board.Closed = trelloBoard.Closed
	board.IDOrganization = trelloBoard.IDOrganization
	board.Pinned = trelloBoard.Pinned
	board.URL = trelloBoard.URL
	board.ShortURL = trelloBoard.ShortURL
	return board
}
