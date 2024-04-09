package models

import (
	"database/sql"
	"log"
	"time"

	"github.com/mstgnz/cronjob/config"
)

type User struct {
	ID        int       `json:"id"`
	Fullname  string    `json:"fullname"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	IsAdmin   bool      `json:"is_admin"`
	LastLogin time.Time `json:"last_login"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
}

type UserLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UserRegister struct {
	Fullname string `json:"fullname" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UserPasswordUpdate struct {
	ID         uint   `json:"id"`
	Password   string `json:"password" validate:"required,min=6"`
	RePassword string `json:"re-password" validate:"required,min=6"`
}

func (u *User) GetUsers() []*User {
	users := []*User{}
	return users
}

func (u *User) CreateUser(register *UserRegister) (*User, error) {
	hashPass := config.HashAndSalt(register.Password)
	rows, err := config.App().DB.Query(config.App().QUERY["USER_INSERT"], register.Fullname, register.Email, hashPass, false)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)
	for rows.Next() {
		if err := rows.Scan(&u); err != nil {
			return nil, err
		}
	}
	return u, nil
}

func (u *User) Exists(email string) bool {
	exists := 0
	rows, err := config.App().DB.Query(config.App().QUERY["USER_EXISTS_WITH_EMAIL"], email)
	if err != nil {
		log.Println("User Exists: ", err)
		return exists > 0
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)
	for rows.Next() {
		if err := rows.Scan(&exists); err != nil {
			log.Println("User Exists Scan: ", err)
		}
	}
	return exists > 0
}

func (u *User) GetUserWithId(id int) *User {
	rows, err := config.App().DB.Query(config.App().QUERY["USER_GET_WITH_ID"], id)
	if err != nil {
		log.Println("GetUserWithId: ", err)
		return nil
	}
	for rows.Next() {
		if err := rows.Scan(&u); err != nil {
			log.Println("User Exists Scan: ", err)
		}
	}
	return u
}

func (u *User) GetUserWithMail(email string) *User {
	rows, err := config.App().DB.Query(config.App().QUERY["USER_GET_WITH_EMAIL"], email)
	if err != nil {
		log.Println("GetUserWithId: ", err)
		return nil
	}
	for rows.Next() {
		if err := rows.Scan(&u); err != nil {
			log.Println("User Exists Scan: ", err)
		}
	}
	return u
}

func (u *User) UpdateUser(fullname, email string) *User {
	return u
}

func (u *User) UpdateUserPassword(password string) *User {
	return u
}

func (u *User) DeleteUser(id int) *User {
	return u
}
