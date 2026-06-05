package models

type User struct {
	ID        int    `db:"id"`
	Username  string `db:"username"`
	Email     string `db:"email"`
	CreatedAt string `db:"created_at"`
	CreatedBy int    `db:"created_by"`
	UpdatedAt string `db:"updated_at"`
	UpdatedBy int    `db:"updated_by"`
	DeletedAt string `db:"deleted_at"`
	Deleted   bool   `db:"deleted"`
	DeletedBy int    `db:"deleted_by"`
	Role      string `db:"role"`
}
