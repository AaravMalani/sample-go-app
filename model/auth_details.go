package model

type AuthDetails struct {
	ID       string `json:"id"`
	Salt     string `json:"salt"`
	Password string `json:"password"`
}
