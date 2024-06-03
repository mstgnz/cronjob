package models

import (
	"log"
	"time"

	"github.com/mstgnz/cronjob/config"
)

type User struct {
	ID        int        `json:"id"`
	Fullname  string     `json:"fullname"`
	Email     string     `json:"email"`
	Password  string     `json:"-"`
	IsAdmin   bool       `json:"is_admin"`
	LastLogin *time.Time `json:"last_login,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
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
	defer func() {
		_ = rows.Close()
	}()
	for rows.Next() {
		if err := rows.Scan(&u.ID, &u.Fullname, &u.Email, &u.Password, &u.IsAdmin, &u.LastLogin, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt); err != nil {
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
	defer func() {
		_ = rows.Close()
	}()
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
	defer func() {
		_ = rows.Close()
	}()
	for rows.Next() {
		if err := rows.Scan(&u.ID, &u.Fullname, &u.Email, &u.IsAdmin, &u.Password); err != nil {
			log.Println("User Scan: ", err)
			return nil
		}
	}
	return u
}

func (u *User) GetUserWithMail(email string) *User {
	rows, err := config.App().DB.Query(config.App().QUERY["USER_GET_WITH_EMAIL"], email)
	if err != nil {
		log.Println("GetUserWithMail: ", err)
		return nil
	}
	defer func() {
		_ = rows.Close()
	}()
	for rows.Next() {
		if err := rows.Scan(&u.ID, &u.Fullname, &u.Email, &u.IsAdmin, &u.Password); err != nil {
			log.Println("User Scan: ", err)
			return nil
		}
	}
	return u
}

func (u *User) Update(query string, params []any) (*User, error) {
	rows, err := config.App().DB.Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()
	for rows.Next() {
		if err := rows.Scan(&u.ID, &u.Fullname, &u.Email, &u.IsAdmin, &u.Password); err != nil {
			return nil, err
		}
	}
	return u, nil
}

func (u *User) UpdatePassword(password string) *User {
	return u
}

func (u *User) UpdateLastLogin() *User {
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")
	rows, _ := config.App().DB.Query(config.App().QUERY["USER_LAST_LOGIN"], formattedTime, u.ID)
	defer func() {
		_ = rows.Close()
	}()
	return u
}

func (u *User) DeleteUser(id int) *User {
	return u
}
