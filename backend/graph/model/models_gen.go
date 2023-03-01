// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

type NewOperation struct {
	Description string  `json:"description"`
	Date        int     `json:"date"`
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
}

type Operation struct {
	ID          string  `json:"id"`
	Description string  `json:"description"`
	Date        int     `json:"date"`
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
	UserID      string  `json:"userId"`
}

type UpdateOperation struct {
	ID          string   `json:"id"`
	Description *string  `json:"description"`
	Date        *int     `json:"date"`
	Amount      *float64 `json:"amount"`
	Category    *string  `json:"category"`
}

type Rol string

const (
	RolAdmin Rol = "ADMIN"
	RolUser  Rol = "USER"
)

var AllRol = []Rol{
	RolAdmin,
	RolUser,
}

func (e Rol) IsValid() bool {
	switch e {
	case RolAdmin, RolUser:
		return true
	}
	return false
}

func (e Rol) String() string {
	return string(e)
}

func (e *Rol) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Rol(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Rol", str)
	}
	return nil
}

func (e Rol) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
