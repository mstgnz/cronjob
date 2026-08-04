package models

import (
	"fmt"
	"time"

	"github.com/mstgnz/cronjob/pkg/auth"
	"github.com/mstgnz/cronjob/pkg/config"
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
	// TokensValidAfter rejects tokens issued before it; set on logout and password change.
	TokensValidAfter *time.Time `json:"-"`
	CreatedAt        *time.Time `json:"created_at,omitempty"`
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
	ID int `json:"id" validate:"omitempty"` // This field is required if the administrator wants to update a user.
	// CurrentPassword is required when users change their own password, so a
	// borrowed session cannot be turned into permanent access.
	CurrentPassword string `json:"current-password" validate:"omitempty"`
	Password        string `json:"password" validate:"required,min=6"`
	RePassword      string `json:"re-password" validate:"required,min=6"`
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

	hashPass := auth.HashAndSalt(register.Password)
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
		if err := rows.Scan(&m.ID, &m.Fullname, &m.Email, &m.Phone, &m.IsAdmin, &m.Active, &m.Password, &m.LastLogin, &m.TokensValidAfter, &m.CreatedAt, &m.UpdatedAt); err != nil {
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
		if err := rows.Scan(&m.ID, &m.Fullname, &m.Email, &m.Phone, &m.IsAdmin, &m.Active, &m.Password, &m.LastLogin, &m.TokensValidAfter, &m.CreatedAt, &m.UpdatedAt); err != nil {
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

// InvalidateTokens makes every token issued for this user before now unusable.
func (m *User) InvalidateTokens() error {
	stmt, err := config.App().DB.Prepare(config.App().QUERY["USER_INVALIDATE_TOKENS"])
	if err != nil {
		return err
	}
	defer func() {
		_ = stmt.Close()
	}()

	// UTC on purpose: the column is timestamp without time zone and is compared
	// against the token's iat claim, which is epoch based. Writing local time here
	// would shift the cut-off by the zone offset and retire every valid token.
	_, err = stmt.Exec(time.Now().UTC(), m.ID)
	return err
}

func (m *User) PasswordUpdate(password string) error {
	stmt, err := config.App().DB.Prepare(config.App().QUERY["USER_UPDATE_PASS"])
	if err != nil {
		return err
	}

	// UTC: this value also lands in tokens_valid_after, which is compared against
	// the token's epoch based iat claim.
	updateAt := time.Now().UTC().Format("2006-01-02 15:04:05")
	hashPass := auth.HashAndSalt(password)
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
