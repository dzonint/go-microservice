package handlers

import (
	"net/http"

	"github.com/dzonint/go-microservice/data"
)

type Users struct{}

func NewUsers() *Users {
	return &Users{}
}

func (u *Users) GetUsers(rw http.ResponseWriter, r *http.Request) {
	users, err := data.GetUsers()
	if err != nil {
		http.Error(rw, "Unable to get users", http.StatusInternalServerError)
		return
	}

	err = users.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshall JSON", http.StatusInternalServerError)
		return
	}
}
