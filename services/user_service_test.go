package services

import (
	"database/sql"
	"example/dtos"
	apierrors "example/errors"
	"example/models"
	"testing"
)

type fakeUserRepository struct {
	users       []models.User
	total       int
	user        *models.User
	createdUser *models.User
	updatedUser *models.User
	err         error
	createArg   models.User
}

func (f *fakeUserRepository) GetAllUsers(limit, offset int) ([]models.User, error) {
	return f.users, f.err
}

func (f *fakeUserRepository) CountUsers() (int, error) {
	return f.total, f.err
}

func (f *fakeUserRepository) GetUserById(id int) (*models.User, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.user, nil
}

func (f *fakeUserRepository) CreateUser(user models.User) (*models.User, error) {
	f.createArg = user
	if f.err != nil {
		return nil, f.err
	}
	return f.createdUser, nil
}

func (f *fakeUserRepository) UpdateUser(id int, user models.User) (*models.User, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.updatedUser, nil
}

func (f *fakeUserRepository) DeleteUser(id int) error {
	return f.err
}

func TestGetAllUsersMapsModelsToDTOs(t *testing.T) {
	repo := &fakeUserRepository{
		users: []models.User{
			{ID: 1, Username: "Alice", Email: "alice@example.com", CreatedAt: "2026-06-05", Role: "admin"},
		},
		total: 4,
	}
	service := NewUserService(repo)

	users, total, resp := service.GetAllUsers(10, 0)
	if resp != nil {
		t.Fatalf("expected success, got %#v", resp)
	}
	if total != 4 {
		t.Fatalf("expected total 4, got %d", total)
	}
	if len(users) != 1 || users[0].Username != "Alice" || users[0].Role != "admin" {
		t.Fatalf("unexpected mapped users: %#v", users)
	}
}

func TestGetUserByIDMapsNoRowsToUserNotFound(t *testing.T) {
	service := NewUserService(&fakeUserRepository{err: sql.ErrNoRows})

	user, resp := service.GetUserByID(99)
	if user != nil {
		t.Fatalf("expected nil user, got %#v", user)
	}
	if resp == nil || resp.Code != apierrors.ErrorNotFound.Code {
		t.Fatalf("expected ErrorNotFound, got %#v", resp)
	}
}

func TestCreateUserDefaultsRole(t *testing.T) {
	repo := &fakeUserRepository{
		createdUser: &models.User{ID: 2, Username: "John", Email: "john@example.com", Role: "user"},
	}
	service := NewUserService(repo)

	user, resp := service.CreateUser(dtos.UserDTO{Username: "John", Email: "john@example.com"})
	if resp != nil {
		t.Fatalf("expected success, got %#v", resp)
	}
	if repo.createArg.Role != "user" {
		t.Fatalf("expected create role user, got %q", repo.createArg.Role)
	}
	if user == nil || user.Role != "user" {
		t.Fatalf("expected returned user role user, got %#v", user)
	}
}

func TestUpdateUserMapsNoRowsToUserNotFound(t *testing.T) {
	service := NewUserService(&fakeUserRepository{err: sql.ErrNoRows})

	user, resp := service.UpdateUser(99, dtos.UserDTO{Username: "Missing", Email: "missing@example.com"})
	if user != nil {
		t.Fatalf("expected nil user, got %#v", user)
	}
	if resp == nil || resp.Code != apierrors.ErrorNotFound.Code {
		t.Fatalf("expected ErrorNotFound, got %#v", resp)
	}
}
