package models

import (
	"fmt"
	"time"

	"github.com/mstgnz/cronjob/config"
)

type Group struct {
	ID        int        `json:"id"`
	UID       int        `json:"uid"`
	UserID    int        `json:"user_id"`
	Name      string     `json:"name"`
	Active    bool       `json:"active"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

func (g *Group) Create() (int64, error) {
	stmt, err := config.App().DB.Prepare(config.App().QUERY["GROUP_INSERT"])
	if err != nil {
		return 0, err
	}

	var uid any
	if g.UID == 0 {
		uid = nil
	} else {
		uid = g.UID
	}

	var lastInsertId int64
	err = stmt.QueryRow(uid, g.UserID, g.Name).Scan(&lastInsertId)
	if err != nil {
		return 0, err
	}

	return lastInsertId, nil
}

func (g *Group) NameExists() (bool, error) {
	exists := 0

	// prepare
	stmt, err := config.App().DB.Prepare(config.App().QUERY["GROUP_NAME_EXISTS_WITH_USER"])
	if err != nil {
		return false, err
	}

	// query
	rows, err := stmt.Query(g.Name, g.UserID)
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

func (g *Group) IDExists(id, userID int) (bool, error) {
	exists := 0

	// prepare
	stmt, err := config.App().DB.Prepare(config.App().QUERY["GROUP_ID_EXISTS_WITH_USER"])
	if err != nil {
		return false, err
	}

	// query
	rows, err := stmt.Query(id, userID)
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

func (u *Group) Update(query string, params []any) error {

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
		return fmt.Errorf("Group not updated")
	}

	return nil
}
