package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/mstgnz/cronjob/config"
)

type RequestHeader struct {
	ID        int        `json:"id"`
	RequestID int        `json:"request_id" validate:"required,number"`
	Key       string     `json:"key" validate:"required"`
	Value     string     `json:"value" validate:"required"`
	Active    bool       `json:"active" validate:"boolean"`
	Request   *Request   `json:"request,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type RequestHeaderBulk struct {
	RequestID int    `json:"request_id" validate:"number"`
	Key       string `json:"key" validate:"required"`
	Value     string `json:"value" validate:"required"`
	Active    bool   `json:"active" validate:"boolean"`
}

type RequestHeaderUpdate struct {
	RequestID int    `json:"request_id" validate:"omitempty,number"`
	Key       string `json:"key" validate:"omitempty"`
	Value     string `json:"value" validate:"omitempty"`
	Active    *bool  `json:"active" validate:"omitnil,boolean"`
}

func (m *RequestHeader) Count(userID int) int {
	rowCount := 0

	// prepare count
	stmt, err := config.App().DB.Prepare(config.App().QUERY["REQUEST_HEADERS_COUNT"])
	if err != nil {
		return rowCount
	}

	// query
	rows, err := stmt.Query(userID)
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

func (m *RequestHeader) Paginate(userID, offset, limit int, search string) []*RequestHeader {
	requestHeaders := []*RequestHeader{}

	// prepare requestHeaders paginate
	stmt, err := config.App().DB.Prepare(config.App().QUERY["REQUEST_HEADERS_PAGINATE"])
	if err != nil {
		return requestHeaders
	}

	// query
	rows, err := stmt.Query("%"+search+"%", userID, offset, limit)
	if err != nil {
		return requestHeaders
	}
	defer func() {
		_ = stmt.Close()
		_ = rows.Close()
	}()
	for rows.Next() {
		requestHeader := &RequestHeader{
			Request: &Request{},
		}

		if err := rows.Scan(&requestHeader.ID, &requestHeader.RequestID, &requestHeader.Key, &requestHeader.Value, &requestHeader.Active, &requestHeader.CreatedAt, &requestHeader.UpdatedAt, &requestHeader.DeletedAt, &requestHeader.Request.Url); err != nil {
			return requestHeaders
		}

		requestHeaders = append(requestHeaders, requestHeader)
	}

	return requestHeaders
}

func (m *RequestHeader) Get(userID, id, requestID int, key string) ([]*RequestHeader, error) {

	query := strings.TrimSuffix(config.App().QUERY["REQUEST_HEADERS"], ";")

	if id > 0 {
		query += fmt.Sprintf(" AND rh.id=%d", id)
	}
	if requestID > 0 {
		query += fmt.Sprintf(" AND rh.request_id=%d", requestID)
	}
	if key != "" {
		query += fmt.Sprintf(" AND rh.key='%s'", key)
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

	requestHeaders := []*RequestHeader{}
	for rows.Next() {
		requestHeader := &RequestHeader{}
		if err := rows.Scan(&requestHeader.ID, &requestHeader.RequestID, &requestHeader.Key, &requestHeader.Value, &requestHeader.Active, &requestHeader.CreatedAt, &requestHeader.UpdatedAt, &requestHeader.DeletedAt); err != nil {
			return nil, err
		}
		requestHeaders = append(requestHeaders, requestHeader)
	}

	return requestHeaders, nil
}

func (m *RequestHeader) Create(exec any) error {
	stmt, err := config.App().DB.RunPrepare(exec, config.App().QUERY["REQUEST_HEADER_INSERT"])
	if err != nil {
		return err
	}

	// user_id,url,method,content,active
	err = stmt.QueryRow(m.RequestID, m.Key, m.Value, m.Active).Scan(&m.ID, &m.RequestID, &m.Key, &m.Value, &m.Active)
	if err != nil {
		return err
	}

	return nil
}

func (m *RequestHeader) HeaderExists(exec any, userID int) (bool, error) {
	exists := 0

	// prepare
	stmt, err := config.App().DB.RunPrepare(exec, config.App().QUERY["REQUEST_HEADER_EXISTS_WITH_USER"])
	if err != nil {
		return false, err
	}

	// query
	rows, err := stmt.Query(m.Key, userID, m.RequestID)
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

func (m *RequestHeader) IDExists(id, userID int) (bool, error) {
	exists := 0

	// prepare
	stmt, err := config.App().DB.Prepare(config.App().QUERY["REQUEST_HEADER_ID_EXISTS_WITH_USER"])
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

func (m *RequestHeader) Update(query string, params []any) error {

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
		return fmt.Errorf("Request Header not updated")
	}

	return nil
}

func (m *RequestHeader) Delete(id, userID int) error {

	stmt, err := config.App().DB.Prepare(config.App().QUERY["REQUEST_HEADER_DELETE"])
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
		return fmt.Errorf("Request Header not deleted")
	}

	return nil
}
