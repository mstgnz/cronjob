package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/mstgnz/cronjob/config"
)

type Request struct {
	ID        int        `json:"id"`
	UserID    int        `json:"user_id" validate:"number"`
	Url       string     `json:"url" validate:"required,url"`
	Method    string     `json:"method" validate:"required,oneof=GET POST PUT PATCH"`
	Content   string     `json:"content" validate:"required,json"`
	Active    bool       `json:"active" validate:"boolean"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type RequestUpdate struct {
	UserID  int    `json:"user_id" validate:"omitempty,number"`
	Url     string `json:"url" validate:"omitempty,url"`
	Method  string `json:"method" validate:"omitempty,oneof=GET POST PUT PATCH"`
	Content string `json:"content" validate:"omitempty,json"`
	Active  *bool  `json:"active" validate:"omitnil,boolean"`
}

type RequestBulk struct {
	UserID         int                  `json:"user_id" validate:"number"`
	Url            string               `json:"url" validate:"required,url"`
	Method         string               `json:"method" validate:"required,oneof=GET POST PUT PATCH"`
	Content        string               `json:"content" validate:"required,json"`
	Active         bool                 `json:"active" validate:"boolean"`
	RequestHeaders []*RequestHeaderBulk `json:"request_headers" validate:"required,nonempty,dive"`
}

func (m *Request) Get(userID, id int, url string) ([]Request, error) {

	query := strings.TrimSuffix(config.App().QUERY["REQUESTS"], ";")

	if id > 0 {
		query += fmt.Sprintf(" AND id=%v", id)
	}
	if url != "" {
		query += fmt.Sprintf(" AND url=%v", url)
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

	var requests []Request
	for rows.Next() {
		var request Request
		if err := rows.Scan(&request.ID, &request.UserID, &request.Url, &request.Method, &request.Content, &request.Active, &request.CreatedAt, &request.UpdatedAt, &request.DeletedAt); err != nil {
			return nil, err
		}
		requests = append(requests, request)
	}

	return requests, nil
}

func (m *Request) Create(exec any) error {
	stmt, err := config.App().DB.RunPrepare(exec, config.App().QUERY["REQUEST_INSERT"])
	if err != nil {
		return err
	}

	// user_id,url,method,content,active
	err = stmt.QueryRow(m.UserID, m.Url, m.Method, m.Content, m.Active).Scan(&m.ID, &m.UserID, &m.Url, &m.Method, &m.Content, &m.Active)
	if err != nil {
		return err
	}

	return nil
}

func (m *Request) UrlExists() (bool, error) {
	exists := 0

	// prepare
	stmt, err := config.App().DB.Prepare(config.App().QUERY["REQUEST_URL_EXISTS_WITH_USER"])
	if err != nil {
		return false, err
	}

	// query
	rows, err := stmt.Query(m.Url, m.UserID)
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

func (m *Request) IDExists(id, userID int) (bool, error) {
	exists := 0

	// prepare
	stmt, err := config.App().DB.Prepare(config.App().QUERY["REQUEST_ID_EXISTS_WITH_USER"])
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

func (m *Request) Update(query string, params []any) error {

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
		return fmt.Errorf("Request not updated")
	}

	return nil
}

func (m *Request) Delete(id, userID int) error {

	stmt, err := config.App().DB.Prepare(config.App().QUERY["REQUEST_DELETE"])
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
		return fmt.Errorf("Request not deleted")
	}

	return nil
}
