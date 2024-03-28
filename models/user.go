package models

import "time"

type User struct {
	ID        int       `json:"id"`
	Fullname  string    `json:"fullname"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Url       string    `json:"url"`
	IsAdmin   bool      `json:"is_admin"`
	LastLogin time.Time `json:"last_login"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
}

func (u *User) GetUsers() []*User {
	users := []*User{}
	return users
}

func (u *User) GetUserWithId(id int) *User {
	return u
}

func (u *User) GetUserWithMail(mail string) *User {
	return u
}

func (u *User) UpdateUser(fullname, mail string) *User {
	return u
}

func (u *User) UpdateUserPassword(password string) *User {
	return u
}

func (u *User) DeleteUser(id int) *User {
	return u
}
