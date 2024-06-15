package models

import (
	"fmt"
	"time"

	"github.com/mstgnz/cronjob/config"
)

type User struct {
	ID        int        `json:"id"`
	Fullname  string     `json:"fullname" validate:"required"`
	Email     string     `json:"email" validate:"required,email"`
	Password  string     `json:"-" validate:"required"`
	Phone     string     `json:"phone" validate:"required,e164"`
	IsAdmin   bool       `json:"is_admin"`
	LastLogin *time.Time `json:"last_login,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type UserUpdate struct {
	Fullname string `json:"fullname" validate:"omitempty"`
	Email    string `json:"email" validate:"omitempty,email"`
	Phone    string `json:"phone" validate:"omitempty,e164"`
}

type UserLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UserRegister struct {
	Fullname string `json:"fullname" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Phone    string `json:"phone"`
}

type UserPasswordUpdate struct {
	Password   string `json:"password" validate:"required,min=6"`
	RePassword string `json:"re-password" validate:"required,min=6"`
}

func (m *User) Users() []*User {
	users := []*User{}
	return users
}

func (m *User) Create(register *UserRegister) error {

	stmt, err := config.App().DB.Prepare(config.App().QUERY["USER_INSERT"])
	if err != nil {
		return err
	}

	hashPass := config.HashAndSalt(register.Password)
	err = stmt.QueryRow(register.Fullname, register.Email, hashPass, register.Phone).Scan(&m.ID, &m.Fullname, &m.Email, &m.Phone)
	if err != nil {
		return err
	}

	return nil
}

func (m *User) Exists(email string) (bool, error) {
	exists := 0

	// prepare
	stmt, err := config.App().DB.Prepare(config.App().QUERY["USER_EXISTS_WITH_EMAIL"])
	if err != nil {
		return false, err
	}

	// query
	rows, err := stmt.Query(email)
	if err != nil {
		return false, err
	}
	defer func() {
		_ = stmt.Close()
		_ = rows.Close()
	}()
	for rows.Next() {
		if err := rows.Scan(&exists); err != nil {
			return false, err
		}
	}
	return exists > 0, nil
}

func (m *User) IDExists(id int) (bool, error) {
	exists := 0

	// prepare
	stmt, err := config.App().DB.Prepare(config.App().QUERY["USER_EXISTS_WITH_ID"])
	if err != nil {
		return false, err
	}

	// query
	rows, err := stmt.Query(id)
	if err != nil {
		return false, err
	}
	defer func() {
		_ = stmt.Close()
		_ = rows.Close()
	}()
	for rows.Next() {
		if err := rows.Scan(&exists); err != nil {
			return false, err
		}
	}
	return exists > 0, nil
}

func (m *User) GetWithId(id int) error {

	stmt, err := config.App().DB.Prepare(config.App().QUERY["USER_GET_WITH_ID"])
	if err != nil {
		return err
	}

	rows, err := stmt.Query(id)
	if err != nil {
		return err
	}
	defer func() {
		_ = stmt.Close()
		_ = rows.Close()
	}()

	found := false
	for rows.Next() {
		if err := rows.Scan(&m.ID, &m.Fullname, &m.Email, &m.IsAdmin, &m.Password); err != nil {
			return err
		}
		found = true
	}

	if !found {
		return fmt.Errorf("User Not Found")
	}

	return nil
}

func (m *User) GetWithMail(email string) error {

	stmt, err := config.App().DB.Prepare(config.App().QUERY["USER_GET_WITH_EMAIL"])
	if err != nil {
		return err
	}

	rows, err := stmt.Query(email)
	if err != nil {
		return err
	}
	defer func() {
		_ = stmt.Close()
		_ = rows.Close()
	}()

	found := false
	for rows.Next() {
		if err := rows.Scan(&m.ID, &m.Fullname, &m.Email, &m.IsAdmin, &m.Password); err != nil {
			return err
		}
		found = true
	}

	if !found {
		return fmt.Errorf("User Not Found")
	}

	return nil
}

func (m *User) Update(query string, params []any) error {

	stmt, err := config.App().DB.Prepare(query)
	if err != nil {
		return err
	}

	result, err := stmt.Exec(params...)
	if err != nil {
		return err
	}
	defer func() {
		_ = stmt.Close()
	}()

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return fmt.Errorf("User not updated")
	}

	return nil
}

func (m *User) UpdatePassword(password string) error {
	stmt, err := config.App().DB.Prepare(config.App().QUERY["USER_UPDATE_PASS"])
	if err != nil {
		return err
	}

	updateAt := time.Now().Format("2006-01-02 15:04:05")
	hashPass := config.HashAndSalt(password)
	result, err := stmt.Exec(hashPass, updateAt, m.ID)
	if err != nil {
		return err
	}
	defer func() {
		_ = stmt.Close()
	}()

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return fmt.Errorf("User password not updated")
	}

	return nil
}

func (m *User) UpdateLastLogin() error {
	lastLogin := time.Now().Format("2006-01-02 15:04:05")

	stmt, err := config.App().DB.Prepare(config.App().QUERY["USER_LAST_LOGIN"])
	if err != nil {
		return err
	}

	result, err := stmt.Exec(lastLogin, m.ID)
	if err != nil {
		return err
	}
	defer func() {
		_ = stmt.Close()
	}()

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return fmt.Errorf("User last login not updated")
	}
	return nil
}

func (m *User) Delete(userID int) error {
	stmt, err := config.App().DB.Prepare(config.App().QUERY["USER_DELETE"])
	if err != nil {
		return err
	}

	deleteAndUpdate := time.Now().Format("2006-01-02 15:04:05")

	result, err := stmt.Exec(deleteAndUpdate, deleteAndUpdate, userID)
	if err != nil {
		return err
	}
	defer func() {
		_ = stmt.Close()
	}()

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return fmt.Errorf("User not deleted")
	}

	return nil
}
