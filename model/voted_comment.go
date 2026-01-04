package model

type VotedComment struct {
	Comment
	Vote int16 `json:"vote"`
}
