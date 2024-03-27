package config

import (
	"bufio"
	"errors"
	"io/fs"
	"os"
	"reflect"
	"strings"
)

// LoadSQLQueries loads SQL queries from a file and populates the QUERY map.
func LoadSQLQueries() (map[string]string, error) {
	query := make(map[string]string)
	file, err := os.Open("query.sql")
	if err != nil {
		return query, err
	}
	defer func(file fs.File) {
		_ = file.Close()
	}(file)
	query, err = parseSQLQueries(file, query)
	return query, err
}

// parseSQLQueries reads the SQL queries from the provided file and populates the QUERY map.
func parseSQLQueries(file fs.File, query map[string]string) (map[string]string, error) {
	scanner := bufio.NewScanner(file)
	var key string
	var queries []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if isSQLQuery(line) || len(key) > 0 {
			if len(key) > 0 {
				queries = append(queries, line)
				if strings.HasSuffix(line, ";") {
					query[key] = strings.Join(queries, " ")
					key, queries = "", nil
				}
			} else {
				key = extractKey(line)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return query, errors.New("error reading file: " + err.Error())
	}
	return query, nil
}

// isSQLQuery checks if the given line is an SQL query or a comment.
func isSQLQuery(line string) bool {
	return HasPrefixInList(line, []string{"-- ", "SELECT", "INSERT", "UPDATE", "DELETE"})
}

// extractKey extracts the key from the comment line.
func extractKey(line string) string {
	if strings.HasPrefix(line, "-- ") {
		return strings.Split(line, "-- ")[1]
	}
	return ""
}

// HasPrefixInList is a prefix checker
func HasPrefixInList(str string, prefixes []string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(str, prefix) {
			return true
		}
	}
	return false
}

// ContainsElement checks if a value exists in a slice.
func ContainsElement(val interface{}, array interface{}) bool {
	arr := reflect.ValueOf(array)
	if arr.Kind() != reflect.Slice {
		return false
	}
	for i := 0; i < arr.Len(); i++ {
		if reflect.DeepEqual(val, arr.Index(i).Interface()) {
			return true
		}
	}
	return false
}
