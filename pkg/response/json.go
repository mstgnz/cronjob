package response

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
)

func ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576 // one megabyte

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)

	if err := dec.Decode(data); err != nil {
		return err
	}

	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return errors.New("body must have only a single JSON value")
	}
	return nil
}

func WriteJSON(w http.ResponseWriter, httpStatus int, data Response, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func ResponseToJSON(response Response) []byte {
	json, err := json.Marshal(response)
	if err != nil {
		return []byte("{status:false, message:response convert error, data:null}")
	}
	return json
}

func MaptoJSON(response map[string]any) []byte {
	json, err := json.Marshal(response)
	if err != nil {
		return []byte("{status:false, message:data convert error, data:null}")
	}
	return json
}

// ConvertStringIDsToInt converts specified string IDs in JSON data to integers
func ConvertStringIDsToInt(r *http.Request, columns ...string) ([]byte, error) {
	var resultJSON []byte
	var data map[string]any

	// Decode body
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return resultJSON, err
	}

	// Convert specified string IDs to integers
	for _, column := range columns {
		if idStr, ok := data[column].(string); ok {
			idInt, err := strconv.Atoi(idStr)
			if err != nil {
				return resultJSON, err
			}
			data[column] = idInt
		}
	}

	// Marshal data back to JSON
	resultJSON, err = json.Marshal(data)
	if err != nil {
		return resultJSON, err
	}

	return resultJSON, nil
}

// ConvertStringBoolsToBool converts specified string Bools in JSON data to bool
func ConvertStringBoolsToBool(r *http.Request, columns ...string) ([]byte, error) {
	var resultJSON []byte
	var data map[string]any

	// Decode body
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return resultJSON, err
	}

	// Convert specified string Bools to bool
	for _, column := range columns {
		if val, ok := data[column].(string); ok {
			boolValue, err := strconv.ParseBool(val)
			if err != nil {
				return resultJSON, err
			}
			data[column] = boolValue
		}
	}

	// Marshal data back to JSON
	resultJSON, err = json.Marshal(data)
	if err != nil {
		return resultJSON, err
	}

	return resultJSON, nil
}
