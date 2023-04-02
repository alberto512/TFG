package model

type User struct {
	ID         string       `json:"id"`
	Username   string       `json:"username"`
	Password   string       `json:"password"`
	Role       Role         `json:"role"`
	Operations []*Operation `json:"operations"`
}