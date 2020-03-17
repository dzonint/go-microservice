package data

import (
	"encoding/json"
	"io"
	"time"

	"github.com/asdine/storm"
)

type User struct {
	ID        int       `json:"id" storm:"id,increment=1"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Gender    string    `json:"gender"`
	IPAddress string    `json:"ip_address"`
	CreatedAt time.Time `json:"created_at"`
}

type Users []User

func (u *Users) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(u)
}

func (u *User) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(u)
}

func AddUser(u *User) error {
	db, err := storm.Open("users.db")
	if err != nil {
		return ErrFailedToOpenDB
	}
	defer db.Close()

	u.CreatedAt = time.Now()

	err = db.Save(u)
	if err != nil {
		return ErrFailedToAddUser
	}

	return nil
}

func GetUsers() (Users, error) {
	var user []User
	db, err := storm.Open("users.db")
	if err != nil {
		return []User{}, ErrFailedToOpenDB
	}
	defer db.Close()

	err = db.All(&user)
	if err != nil {
		return []User{}, ErrFailedToGetUsers
	}

	return user, nil
}
