package services

import (
	"database/sql"
	"errors"
	"example/dtos"
	apierrors "example/errors"
	"example/models"
	"example/repositories"
)

type UserService interface {
	GetAllUsers(limit, offset int) ([]dtos.UserDTO, int, *apierrors.CommonResponse)
	GetUserByID(id int) (*dtos.UserDTO, *apierrors.CommonResponse)
	CreateUser(userDTO dtos.UserDTO) (*dtos.UserDTO, *apierrors.CommonResponse)
	UpdateUser(id int, userDTO dtos.UserDTO) (*dtos.UserDTO, *apierrors.CommonResponse)
	DeleteUser(id int) *apierrors.CommonResponse
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetAllUsers(limit, offset int) ([]dtos.UserDTO, int,  *apierrors.CommonResponse ) {
	users, err := s.repo.GetAllUsers(limit, offset)
	if err != nil {
		return nil, 0, &apierrors.ErrorInternal
	}

	total, err := s.repo.CountUsers()
	if err != nil {	
		return nil, 0, &apierrors.ErrorInternal
	}

	var userDTOs []dtos.UserDTO
	for _, user := range users {
		userDTOs = append(userDTOs, dtos.UserDTO{
			ID: user.ID, Username: user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			Role:      user.Role,
		})
	}

	return userDTOs, total, nil
}

func (s *userService) GetUserByID(id int) (*dtos.UserDTO, *apierrors.CommonResponse) {
	// Implementation for fetching a user by ID
	user, err := s.repo.GetUserById(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &apierrors.ErrorNotFound
		}
		return nil, &apierrors.ErrorInternal
	}

	userDTO := &dtos.UserDTO{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		Role:      user.Role,
	}

	return userDTO, nil

}

func (s *userService) CreateUser(userDTO dtos.UserDTO) (*dtos.UserDTO, *apierrors.CommonResponse) {
	// Implementation for creating a new user
	if userDTO.Role == "" {
		userDTO.Role = "user" // Default role if not provided
	}

	newUser, err := s.repo.CreateUser(models.User{
		Username:  userDTO.Username,
		Email:     userDTO.Email,
		CreatedAt: userDTO.CreatedAt,
		Role:      userDTO.Role,
	})
	if err != nil {
		return nil, &apierrors.ErrorInternal
	}

	return &dtos.UserDTO{
		ID:        newUser.ID,
		Username:  newUser.Username,
		Email:     newUser.Email,
		CreatedAt: newUser.CreatedAt,
		Role:      newUser.Role,
	}, nil
}

func (s *userService) UpdateUser(id int, userDTO dtos.UserDTO) (*dtos.UserDTO, *apierrors.CommonResponse) {
	// Implementation for updating an existing user
	updatedUser, err := s.repo.UpdateUser(id, models.User{
		Username:  userDTO.Username,
		Email:     userDTO.Email,
		CreatedAt: userDTO.CreatedAt,
		Role:      userDTO.Role,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &apierrors.ErrorNotFound
		}
		return nil, &apierrors.ErrorInternal
	}

	return &dtos.UserDTO{
		ID:        updatedUser.ID,
		Username:  updatedUser.Username,
		Email:     updatedUser.Email,
		CreatedAt: updatedUser.CreatedAt,
		Role:      updatedUser.Role,
	}, nil
}

func (s *userService) DeleteUser(id int) *apierrors.CommonResponse {
	// Implementation for deleting a user
	if err := s.repo.DeleteUser(id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &apierrors.ErrorNotFound
		}
		return &apierrors.ErrorInternal
	}
	return nil
}
