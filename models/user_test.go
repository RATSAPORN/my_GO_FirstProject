package models

import (
	"reflect"
	"testing"
)

func TestUserDBTags(t *testing.T) {
	tests := map[string]string{
		"ID":        "id",
		"Username":  "username",
		"Email":     "email",
		"CreatedAt": "created_at",
		"Role":      "role",
	}

	modelType := reflect.TypeOf(User{})
	for fieldName, want := range tests {
		field, ok := modelType.FieldByName(fieldName)
		if !ok {
			t.Fatalf("field %s not found", fieldName)
		}
		if got := field.Tag.Get("db"); got != want {
			t.Fatalf("field %s db tag = %q, want %q", fieldName, got, want)
		}
	}
}
