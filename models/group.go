package models

import (
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

func (g *Group) Exists() (bool, error) {
	exists := 0

	// prepare
	stmt, err := config.App().DB.Prepare(config.App().QUERY["GROUP_EXISTS_WITH_USER"])
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
