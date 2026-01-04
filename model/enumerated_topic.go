package model

type EnumeratedTopic struct {
	Topic
	Count int `json:"count"`
}
