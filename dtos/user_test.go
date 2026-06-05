package dtos

import (
	"reflect"
	"testing"
)

func TestUserDTOJSONTags(t *testing.T) {
	tests := map[string]string{
		"ID":        "id",
		"Username":  "username",
		"Email":     "email",
		"CreatedAt": "created_at",
		"Role":      "role",
	}

	dtoType := reflect.TypeOf(UserDTO{})
	for fieldName, want := range tests {
		field, ok := dtoType.FieldByName(fieldName)
		if !ok {
			t.Fatalf("field %s not found", fieldName)
		}
		if got := field.Tag.Get("json"); got != want {
			t.Fatalf("field %s json tag = %q, want %q", fieldName, got, want)
		}
	}
}

func TestCommonPaginationDataJSONTags(t *testing.T) {
	tests := map[string]string{
		"Data":   "data",
		"Total":  "total",
		"Limit":  "limit",
		"Offset": "offset",
	}

	dtoType := reflect.TypeOf(CommonPaginationData{})
	for fieldName, want := range tests {
		field, ok := dtoType.FieldByName(fieldName)
		if !ok {
			t.Fatalf("field %s not found", fieldName)
		}
		if got := field.Tag.Get("json"); got != want {
			t.Fatalf("field %s json tag = %q, want %q", fieldName, got, want)
		}
	}
}
