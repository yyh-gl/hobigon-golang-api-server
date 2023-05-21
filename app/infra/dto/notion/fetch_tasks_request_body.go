package notion

// FIXME: Trello -> Notion への移行を突貫工事で作ったのでリファクタ推奨

type FetchTasksRequestBody struct {
	PageSize int    `json:"page_size"`
	Filter   any    `json:"filter,omitempty"`
	Sorts    []Sort `json:"sorts,omitempty"`
}

type SingleFilter struct {
	Property string  `json:"property"`
	Date     *Date   `json:"date,omitempty"`
	Select   *Select `json:"select,omitempty"`
}

type AndFilter struct {
	And any `json:"and"`
}

type OrFilter struct {
	Or any `json:"or"`
}

type Date struct {
	OnOrBefore string `json:"on_or_before,omitempty"`
	OnOrAfter  string `json:"on_or_after,omitempty"`
}

type Select struct {
	Equals string `json:"equals,omitempty"`
}

type Sort struct {
	Property  string `json:"property"`
	Direction string `json:"direction"`
}
