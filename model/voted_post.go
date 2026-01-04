package model

type VotedPost struct {
	Post
	Vote int16 `json:"vote"`
}
