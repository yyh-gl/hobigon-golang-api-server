package notion

// FIXME: Trello -> Notion への移行を突貫工事で作ったのでリファクタ推奨

type FetchTasksRequestBody struct {
	PageSize int    `json:"page_size"`
	Filter   Filter `json:"filter,omitempty"`
	Sorts    []Sort `json:"sorts,omitempty"`
}

type Filter interface{}

type SingleFilter struct {
	Property string `json:"property"`
	Date     Date   `json:"date,omitempty"`
}

type AndFilter struct {
	And []SingleFilter `json:"and"`
}

type Date struct {
	OnOrBefore string `json:"on_or_before,omitempty"`
	OnOrAfter  string `json:"on_or_after,omitempty"`
}

type Sort struct {
	Property  string `json:"property"`
	Direction string `json:"direction"`
}
