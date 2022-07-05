package domain

import "time"

type AnswerRequest struct {
	// the id for create a new event
	ID string `json:"id,omitempty"`

	// data is answers list
	Data map[string]string `json:"data" binding:"required"`
}

type Response struct {
	// event id
	ID string `json:"answer_id"`
	// type of event
	Event string `json:"event"`
	// data is answers list
	Data     map[string]string `json:"data" example:"en:Map,ru:Карта,kk:Карталар"`
	CreateAt time.Time         `json:"createdAt,omitempty"`
	Code     string
	Content  string
}
