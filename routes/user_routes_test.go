package routes

import (
	"bytes"
	"encoding/json"
	"example/controller"
	"example/dtos"
	apierrors "example/errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type fakeUserService struct{}

func (f *fakeUserService) GetAllUsers(limit, offset int) ([]dtos.UserDTO, int, *apierrors.CommonResponse) {
	return []dtos.UserDTO{{ID: 1, Username: "Alice", Email: "alice@example.com", Role: "admin"}}, 1, nil
}

func (f *fakeUserService) GetUserByID(id int) (*dtos.UserDTO, *apierrors.CommonResponse) {
	return &dtos.UserDTO{ID: id, Username: "Alice", Email: "alice@example.com", Role: "admin"}, nil
}

func (f *fakeUserService) CreateUser(userDTO dtos.UserDTO) (*dtos.UserDTO, *apierrors.CommonResponse) {
	return &dtos.UserDTO{ID: 2, Username: userDTO.Username, Email: userDTO.Email, Role: "user"}, nil
}

func (f *fakeUserService) UpdateUser(id int, userDTO dtos.UserDTO) (*dtos.UserDTO, *apierrors.CommonResponse) {
	return &dtos.UserDTO{ID: id, Username: userDTO.Username, Email: userDTO.Email, Role: userDTO.Role}, nil
}

func (f *fakeUserService) DeleteUser(id int) *apierrors.CommonResponse {
	return nil
}

func setupTestServer(t *testing.T) *httptest.Server {
	t.Helper()

	gin.SetMode(gin.TestMode)

	service := &fakeUserService{}
	userController := controller.NewUserController(service)
	router := gin.New()
	api := router.Group("/api/v1")
	RegisterUserRoutes(api, userController, service)

	return httptest.NewServer(router)
}

func doJSONRequest(t *testing.T, method, url, body string) (*http.Response, map[string]any) {
	t.Helper()

	req, err := http.NewRequest(method, url, bytes.NewBufferString(body))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("X-User-Id", "1")
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	t.Cleanup(func() {
		resp.Body.Close()
	})

	var responseBody map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	return resp, responseBody
}

func TestRegisterUserRoutesRealAPIRequests(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	tests := []struct {
		name   string
		method string
		path   string
		body   string
	}{
		{
			name:   "get users",
			method: http.MethodGet,
			path:   "/api/v1/users/?page=1&per_page=10",
		},
		{
			name:   "get user by id",
			method: http.MethodGet,
			path:   "/api/v1/users/1",
		},
		{
			name:   "create user",
			method: http.MethodPost,
			path:   "/api/v1/users/",
			body:   `{"username":"Jane","email":"jane@example.com"}`,
		},
		{
			name:   "update user",
			method: http.MethodPut,
			path:   "/api/v1/users/1",
			body:   `{"username":"Alice","email":"alice.updated@example.com","role":"admin"}`,
		},
		{
			name:   "delete user",
			method: http.MethodDelete,
			path:   "/api/v1/users/1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, body := doJSONRequest(t, tt.method, server.URL+tt.path, tt.body)
			t.Logf("%s %s -> status=%d body=%#v", tt.method, tt.path, resp.StatusCode, body)

			if resp.StatusCode != http.StatusOK {
				t.Fatalf("expected status %d, got %d with body %#v", http.StatusOK, resp.StatusCode, body)
			}
			if body["code"] != "0000" || body["message"] != "ok" {
				t.Fatalf("expected success response wrapper, got %#v", body)
			}
		})
	}
}
