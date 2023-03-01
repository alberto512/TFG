package model

type User struct {
	ID         string       `json:"id"`
	Username   string       `json:"username"`
	Password   string       `json:"password"`
	Rol        Rol          `json:"rol"`
	Operations []*Operation `json:"operations"`
}