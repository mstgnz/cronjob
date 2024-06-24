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
	Active    bool       `json:"active"`
	IsAdmin   bool       `json:"is_admin"`
	LastLogin *time.Time `json:"last_login,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type Login struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type Register struct {
	Fullname string `json:"fullname" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Phone    string `json:"phone" validate:"required,e164"`
}

type ProfileUpdate struct {
	ID       int    `json:"id" validate:"omitempty"` // This field is required if the administrator wants to update a user.
	Fullname string `json:"fullname" validate:"omitempty"`
	Email    string `json:"email" validate:"omitempty,email"`
	Phone    string `json:"phone" validate:"omitempty,e164"`
}

type PasswordUpdate struct {
	ID         int    `json:"id" validate:"omitempty"` // This field is required if the administrator wants to update a user.
	Password   string `json:"password" validate:"required,min=6"`
	RePassword string `json:"re-password" validate:"required,min=6"`
}

func (m *User) Count() int {
	rowCount := 0

	// prepare count
	stmt, err := config.App().DB.Prepare(config.App().QUERY["USERS_COUNT"])
	if err != nil {
		return rowCount
	}

	// query
	rows, err := stmt.Query()
	if err != nil {
		return rowCount
	}
	defer func() {
		_ = stmt.Close()
		_ = rows.Close()
	}()
	for rows.Next() {
		if err := rows.Scan(&rowCount); err != nil {
			return rowCount
		}
	}

	return rowCount
}

func (m *User) Get(offset, limit int, search string) []*User {
	users := []*User{}

	// prepare users paginate
	stmt, err := config.App().DB.Prepare(config.App().QUERY["USERS_PAGINATE"])
	if err != nil {
		return users
	}

	// query
	rows, err := stmt.Query("%"+search+"%", offset, limit)
	if err != nil {
		return users
	}
	defer func() {
		_ = stmt.Close()
		_ = rows.Close()
	}()
	for rows.Next() {
		user := &User{}
		if err := rows.Scan(&user.ID, &user.Fullname, &user.Email, &user.Password, &user.Phone, &user.IsAdmin, &user.Active, &user.LastLogin, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt); err != nil {
			return users
		}
		users = append(users, user)
	}

	return users
}

func (m *User) Create(register *Register) error {

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

func (m *User) ProfileUpdate(query string, params []any) error {

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

func (m *User) PasswordUpdate(password string) error {
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

func (m *User) LastLoginUpdate() error {
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

	result, err := stmt.Exec(false, deleteAndUpdate, deleteAndUpdate, userID)
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
