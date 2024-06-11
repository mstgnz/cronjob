package models

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/mstgnz/cronjob/config"
)

type Group struct {
	ID        int        `json:"id"`
	UID       int        `json:"uid" validate:"number"`
	UserID    int        `json:"user_id" validate:"number"`
	Name      string     `json:"name" validate:"required"`
	Active    bool       `json:"active" validate:"required,boolean"`
	Parent    *Group     `json:"parent,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type GroupUpdate struct {
	UID    int    `json:"uid" validate:"omitempty,number"`
	UserID int    `json:"user_id" validate:"omitempty,number"`
	Name   string `json:"name" validate:"omitempty"`
	Active *bool  `json:"active" validate:"omitnil,boolean"`
}

func (m *Group) Get(userID, id, uid int) ([]Group, error) {

	query := strings.TrimSuffix(config.App().QUERY["GROUPS"], ";")

	if id > 0 {
		query += fmt.Sprintf(" AND id=%v", id)
	}
	if uid > 0 {
		query += fmt.Sprintf(" AND uid=%v", uid)
	}

	// prepare
	stmt, err := config.App().DB.Prepare(query)
	if err != nil {
		return nil, err
	}

	// query
	rows, err := stmt.Query(userID)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = stmt.Close()
		_ = rows.Close()
	}()

	var groups []Group
	for rows.Next() {
		var group Group
		var parentID sql.NullInt64
		if err := rows.Scan(&group.ID, &parentID, &group.Name, &group.Active, &group.CreatedAt, &group.UpdatedAt); err != nil {
			return nil, err
		}
		if parentID.Valid {
			group.UID = int(parentID.Int64)
			var parent Group
			row := config.App().DB.QueryRow(config.App().QUERY["GROUPS_WITH_ID"], userID, parentID.Int64)
			if row.Err() != nil {
				return nil, err
			}
			row.Scan(&parent.ID, &parentID, &parent.Name, &parent.Active, &parent.CreatedAt, &parent.UpdatedAt)
			group.Parent = &parent
		}
		groups = append(groups, group)
	}

	return groups, nil
}

func (m *Group) Create() error {
	stmt, err := config.App().DB.Prepare(config.App().QUERY["GROUP_INSERT"])
	if err != nil {
		return err
	}

	var uid any
	if m.UID == 0 {
		uid = nil
	} else {
		uid = m.UID
	}

	var parentID sql.NullInt64
	err = stmt.QueryRow(uid, m.UserID, m.Name).Scan(&m.ID, &parentID, &m.Name, &m.Active)
	if err != nil {
		return err
	}
	if parentID.Valid {
		m.UID = int(parentID.Int64)
	}

	return nil
}

func (m *Group) NameExists() (bool, error) {
	exists := 0

	// prepare
	stmt, err := config.App().DB.Prepare(config.App().QUERY["GROUP_NAME_EXISTS_WITH_USER"])
	if err != nil {
		return false, err
	}

	// query
	rows, err := stmt.Query(m.Name, m.UserID)
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

func (m *Group) IDExists(id, userID int) (bool, error) {
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

func (m *Group) Update(query string, params []any) error {

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

func (m *Group) Delete(id, userID int) error {

	stmt, err := config.App().DB.Prepare(config.App().QUERY["GROUP_DELETE"])
	if err != nil {
		return err
	}

	deleteAndUpdate := time.Now().Format("2006-01-02 15:04:05")

	result, err := stmt.Exec(deleteAndUpdate, deleteAndUpdate, id, userID)
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
		return fmt.Errorf("Group not deleted")
	}

	return nil
}
