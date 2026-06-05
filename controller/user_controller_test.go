package controller

import (
	"bytes"
	"encoding/json"
	"example/dtos"
	apierrors "example/errors"
	"example/services"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type fakeUserService struct {
	users       []dtos.UserDTO
	total       int
	user        *dtos.UserDTO
	createdUser *dtos.UserDTO
	updatedUser *dtos.UserDTO
	resp        apierrors.CommonResponse
	limit       int
	offset      int
	deletedID   int
}

// respOrSuccess lets tests leave resp unset (zero value) to mean a successful call.
func (f *fakeUserService) respOrSuccess() apierrors.CommonResponse {
	if f.resp == (apierrors.CommonResponse{}) {
		return apierrors.SuccessResponse
	}
	return f.resp
}

func (f *fakeUserService) GetAllUsers(limit, offset int) ([]dtos.UserDTO, int, apierrors.CommonResponse) {
	f.limit = limit
	f.offset = offset
	return f.users, f.total, f.respOrSuccess()
}

func (f *fakeUserService) GetUserByID(id int) (*dtos.UserDTO, apierrors.CommonResponse) {
	if resp := f.respOrSuccess(); resp != apierrors.SuccessResponse {
		return nil, resp
	}
	return f.user, apierrors.SuccessResponse
}

func (f *fakeUserService) CreateUser(userDTO dtos.UserDTO) (*dtos.UserDTO, apierrors.CommonResponse) {
	if resp := f.respOrSuccess(); resp != apierrors.SuccessResponse {
		return nil, resp
	}
	return f.createdUser, apierrors.SuccessResponse
}

func (f *fakeUserService) UpdateUser(id int, userDTO dtos.UserDTO) (*dtos.UserDTO, apierrors.CommonResponse) {
	if resp := f.respOrSuccess(); resp != apierrors.SuccessResponse {
		return nil, resp
	}
	return f.updatedUser, apierrors.SuccessResponse
}

func (f *fakeUserService) DeleteUser(id int) apierrors.CommonResponse {
	f.deletedID = id
	return f.respOrSuccess()
}

func setupUserControllerRouter(service services.UserService) *gin.Engine {
	gin.SetMode(gin.TestMode)

	controller := NewUserController(service)
	router := gin.New()
	router.GET("/users/", controller.GetAllUsers)
	router.GET("/users/:id", controller.GetUserByID)
	router.POST("/users/", controller.CreateUser)
	router.PUT("/users/:id", controller.UpdateUser)
	router.DELETE("/users/:id", controller.DeleteUser)
	return router
}

func decodeResponse(t *testing.T, recorder *httptest.ResponseRecorder) map[string]any {
	t.Helper()

	var body map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	return body
}

func TestGetAllUsersSuccess(t *testing.T) {
	service := &fakeUserService{
		users: []dtos.UserDTO{
			{ID: 1, Username: "Alice", Email: "alice@example.com", Role: "admin"},
		},
		total: 12,
	}
	router := setupUserControllerRouter(service)

	req := httptest.NewRequest(http.MethodGet, "/users/?page=2&per_page=5", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
	if service.limit != 5 || service.offset != 5 {
		t.Fatalf("expected limit=5 offset=5, got limit=%d offset=%d", service.limit, service.offset)
	}

	body := decodeResponse(t, recorder)
	if body["code"] != "0000" || body["message"] != "ok" {
		t.Fatalf("expected success wrapper, got %#v", body)
	}

	data := body["data"].(map[string]any)
	if data["total"] != float64(12) || data["limit"] != float64(5) || data["offset"] != float64(5) {
		t.Fatalf("unexpected pagination data: %#v", data)
	}
}

func TestGetAllUsersInvalidPage(t *testing.T) {
	router := setupUserControllerRouter(&fakeUserService{})

	req := httptest.NewRequest(http.MethodGet, "/users/?page=0&per_page=5", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}

	body := decodeResponse(t, recorder)
	if body["code"] != "E100" || body["message"] != "invalid page number" {
		t.Fatalf("expected bad request wrapper, got %#v", body)
	}
}

func TestGetUserByIDNotFound(t *testing.T) {
	router := setupUserControllerRouter(&fakeUserService{resp: apierrors.ErrorNotFound})

	req := httptest.NewRequest(http.MethodGet, "/users/99", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, recorder.Code)
	}

	body := decodeResponse(t, recorder)
	if body["code"] != "E101" || body["message"] != "resource not found" {
		t.Fatalf("expected not found wrapper, got %#v", body)
	}
}

func TestCreateUserSuccess(t *testing.T) {
	service := &fakeUserService{
		createdUser: &dtos.UserDTO{ID: 3, Username: "Jane", Email: "jane@example.com", Role: "user"},
	}
	router := setupUserControllerRouter(service)

	body := bytes.NewBufferString(`{"username":"Jane","email":"jane@example.com"}`)
	req := httptest.NewRequest(http.MethodPost, "/users/", body)
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	resp := decodeResponse(t, recorder)
	if resp["code"] != "0000" || resp["message"] != "ok" {
		t.Fatalf("expected success wrapper, got %#v", resp)
	}
}

func TestUpdateUserSuccess(t *testing.T) {
	service := &fakeUserService{
		updatedUser: &dtos.UserDTO{ID: 2, Username: "John", Email: "john.updated@example.com", Role: "user"},
	}
	router := setupUserControllerRouter(service)

	body := bytes.NewBufferString(`{"username":"John","email":"john.updated@example.com","role":"user"}`)
	req := httptest.NewRequest(http.MethodPut, "/users/2", body)
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	resp := decodeResponse(t, recorder)
	if resp["code"] != "0000" || resp["message"] != "ok" {
		t.Fatalf("expected success wrapper, got %#v", resp)
	}
}

func TestDeleteUserSuccess(t *testing.T) {
	service := &fakeUserService{}
	router := setupUserControllerRouter(service)

	req := httptest.NewRequest(http.MethodDelete, "/users/2", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
	if service.deletedID != 2 {
		t.Fatalf("expected deleted id 2, got %d", service.deletedID)
	}

	resp := decodeResponse(t, recorder)
	if resp["code"] != "0000" || resp["message"] != "ok" {
		t.Fatalf("expected success wrapper, got %#v", resp)
	}
}
