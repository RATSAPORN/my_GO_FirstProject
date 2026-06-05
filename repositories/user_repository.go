package repositories

import (
	"database/sql"
	"example/models"

	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	GetAllUsers(limit, offset int) ([]models.User, error)
	CountUsers() (int, error)
	GetUserById(id int) (*models.User, error)
	CreateUser(user models.User) (*models.User, error)
	UpdateUser(id int, user models.User) (*models.User, error)
	DeleteUser(id int) error
}

type userRepo struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) CountUsers() (int, error) {
	var count int
	err := r.db.Get(&count, "SELECT COUNT(*) FROM users WHERE deleted = false")
	if err != nil {
		return 0, err
	}
	return count, err
}

func (r *userRepo) GetAllUsers(limit, offset int) ([]models.User, error) {
	var users []models.User
	err := r.db.Select(&users, "SELECT id, username, email, created_at, role FROM users WHERE deleted = false ORDER BY id LIMIT :limit OFFSET :offset", map[string]any{"limit": limit, "offset": offset})
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepo) GetUserById(id int) (*models.User, error) {
	// sqlx.Named expands :id and gives positional args; Rebind renumbers to $1 for Postgres.
	query, args, err := sqlx.Named(
		"SELECT id, username, email, created_at, role FROM users WHERE id = :id AND deleted = false",
		map[string]any{"id": id},
	)
	if err != nil {
		return nil, err
	}
	query = r.db.Rebind(query)

	var user models.User
	if err := r.db.Get(&user, query, args...); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) CreateUser(user models.User) (*models.User, error) {
	// NamedQuery binds :username, :email, ... from the struct's db tags and returns the new row.
	rows, err := r.db.NamedQuery(
		`INSERT INTO users (username, email, created_at, created_by, role)
		 VALUES (:username, :email, :created_at, :created_by, :role)
		 RETURNING id, username, email, created_at, role`,
		user,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var newUser models.User
	if rows.Next() {
		if err := rows.StructScan(&newUser); err != nil {
			return nil, err
		}
	}
	return &newUser, nil
}

func (r *userRepo) UpdateUser(id int, user models.User) (*models.User, error) {
	user.ID = id // bind :id from the struct
	rows, err := r.db.NamedQuery(
		`UPDATE users
		 SET username = :username, email = :email, updated_at = CURRENT_TIMESTAMP
		 WHERE id = :id AND deleted = false
		 RETURNING id, username, email, created_at, role`,
		user,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var updatedUser models.User
	if rows.Next() {
		if err := rows.StructScan(&updatedUser); err != nil {
			return nil, err
		}
	} else {
		return nil, sql.ErrNoRows
	}
	return &updatedUser, nil
}

func (r *userRepo) DeleteUser(id int) error {
	result, err := r.db.NamedExec(
		"UPDATE users SET deleted = true, deleted_at = NOW(), deleted_by = :deleted_by WHERE id = :id AND deleted = false",
		map[string]any{"deleted_by": 1, "id": id},
	)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows // or services.ErrUserNot
	}
	return nil
}
