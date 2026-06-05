package middleware

import (
	"encoding/json"
	"example/dtos"
	"example/services"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type fakeUserService struct {
	user *dtos.UserDTO
	err  error
}

func (f *fakeUserService) GetAllUsers(limit, offset int) ([]dtos.UserDTO, int, error) {
	return nil, 0, f.err
}

func (f *fakeUserService) GetUserByID(id int) (*dtos.UserDTO, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.user, nil
}

func (f *fakeUserService) CreateUser(userDTO dtos.UserDTO) (*dtos.UserDTO, error) {
	return nil, f.err
}

func (f *fakeUserService) UpdateUser(id int, userDTO dtos.UserDTO) (*dtos.UserDTO, error) {
	return nil, f.err
}

func (f *fakeUserService) DeleteUser(id int) error {
	return f.err
}

func setupMiddlewareRouter(service services.UserService, handlers ...gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	allHandlers := append([]gin.HandlerFunc{AuthMiddleware(service)}, handlers...)
	allHandlers = append(allHandlers, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	router.GET("/users/:id", allHandlers...)
	return router
}

func decodeMiddlewareResponse(t *testing.T, recorder *httptest.ResponseRecorder) map[string]any {
	t.Helper()

	var body map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	return body
}

func TestAuthMiddlewareRequiresUserIDHeader(t *testing.T) {
	router := setupMiddlewareRouter(&fakeUserService{})

	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, recorder.Code)
	}
	body := decodeMiddlewareResponse(t, recorder)
	if body["code"] != "E103" || body["message"] != "X-User-Id header is required" {
		t.Fatalf("unexpected response: %#v", body)
	}
}

func TestRequireAdminAllowsAdmin(t *testing.T) {
	router := setupMiddlewareRouter(
		&fakeUserService{user: &dtos.UserDTO{ID: 1, Role: "admin"}},
		RequireAdmin(),
	)

	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	req.Header.Set("X-User-Id", "1")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
}

func TestRequireAdminRejectsUser(t *testing.T) {
	router := setupMiddlewareRouter(
		&fakeUserService{user: &dtos.UserDTO{ID: 2, Role: "user"}},
		RequireAdmin(),
	)

	req := httptest.NewRequest(http.MethodGet, "/users/2", nil)
	req.Header.Set("X-User-Id", "2")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, recorder.Code)
	}
}

func TestRequireSelfOrAdminRejectsOtherUser(t *testing.T) {
	router := setupMiddlewareRouter(
		&fakeUserService{user: &dtos.UserDTO{ID: 2, Role: "user"}},
		RequireSelfOrAdmin(),
	)

	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	req.Header.Set("X-User-Id", "2")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, recorder.Code)
	}
}
