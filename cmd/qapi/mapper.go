package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func NewMap() *Map {
	return &Map{}
}

type Map struct {
	Entries []MapEntry
}

type MapEntry struct {
	Required bool
	Field    string
	Dst      interface{}
}

func (m *Map) Required(field string, dst interface{}) {
	m.Entries = append(m.Entries, MapEntry{Required: true, Field: field, Dst: dst})
}

func (m *Map) Optional(field string, dst interface{}) {
	m.Entries = append(m.Entries, MapEntry{Required: false, Field: field, Dst: dst})
}

func JSONMap(r *http.Request, m *Map) error {
	expected, got := "application/json", r.Header.Get("Content-Type")
	if got != expected {
		return fmt.Errorf("Bad content type. Expected: %s, but got: %s", expected, got)
	}
	const maxSize = 1024 * 1024
	d := json.NewDecoder(io.LimitReader(r.Body, maxSize))
	data := map[string]interface{}{}
	if err := d.Decode(&data); err != nil {
		return err
	}
	return GenericMap(data, m)
}

func GenericMap(data map[string]interface{}, m *Map) error {
	for dataField, _ := range data {
		found := false
		for _, entry := range m.Entries {
			if dataField == entry.Field {
				found = true
			}
		}
		if !found {
			return fmt.Errorf("Unknown field: %s", dataField)
		}
	}
	for _, entry := range m.Entries {
		val, ok := data[entry.Field]
		if !ok {
			if entry.Required {
				return fmt.Errorf("Required field is missing: %s", entry.Field)
			}
			continue
		}
		switch dst := entry.Dst.(type) {
		case *string:
			valT, ok := val.(string)
			if !ok {
				return fmt.Errorf("Expected string, but got %T for field: %s", val, entry.Field)
			}
			if valT == "" && entry.Required {
				return fmt.Errorf("Required field is empty: %s", entry.Field)
			}
			*dst = valT
		case *[]string:
			// @TODO Support for []string{} (not needed for now b/c encoding/json
			// doesn't use it)
			valT, ok := val.([]interface{})
			if !ok {
				return fmt.Errorf("Expected []string, but got %T for field: %s", val, entry.Field)
			}
			if len(valT) == 0 && entry.Required {
				return fmt.Errorf("Required field is empty: %s", entry.Field)
			}
			*dst = nil
			for _, el := range valT {
				elT, ok := el.(string)
				if !ok {
					return fmt.Errorf("Expected []string, but got %T for field: %s", el, entry.Field)
				}
				*dst = append(*dst, elT)
			}
		case *time.Time:
		case **time.Time:
			valT, ok := val.(string)
			if !ok {
				return fmt.Errorf("Expected RFC3339 string, but got %T for field: %s", val, entry.Field)
			}
			valTT, err := time.Parse(time.RFC3339, valT)
			if err != nil {
				return fmt.Errorf("Expected RFC3339 string, but got %s for field: %s", err, entry.Field)
			}
			*dst = &valTT
		default:
			return fmt.Errorf("Unsupported map type: %T", entry.Dst)
		}
	}
	return nil
}
