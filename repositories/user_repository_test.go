package repositories

import "testing"

func TestNewUserRepositoryReturnsRepository(t *testing.T) {
	repo := NewUserRepository(nil)
	if repo == nil {
		t.Fatal("expected repository, got nil")
	}
}
