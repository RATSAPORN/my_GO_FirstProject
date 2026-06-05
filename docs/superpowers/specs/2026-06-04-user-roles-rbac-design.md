# User Roles & RBAC — Design

**Date:** 2026-06-04
**Status:** Approved (pending spec review)

## Goal

Make the users API role-based. Each user carries a role (`admin` or `user`), and
middleware enforces what each role may do. Roles are stored in a lookup table and
referenced by foreign key from `users`.

## Decisions

| Topic | Decision |
|-------|----------|
| Scope | Role column **plus** middleware enforcement |
| Roles | `admin`, `user`; default `user` |
| Storage | Lookup table `roles` + FK from `users.role` |
| Identity | `X-User-Id` request header (dev-grade, no auth system) |
| Policy | admin: anything; user: read/update **own** record only |

`X-User-Id` is explicitly a development-grade identity mechanism: it is spoofable and
is not a substitute for real authentication. It exists to wire and test the
enforcement layer. Replacing it with JWT auth later only changes `AuthMiddleware`.

## Data model

`roles` is a lookup table keyed by name. `users.role` is a foreign key to `roles.name`.
A FK to a unique text column is valid in Postgres; this allows a plain string `DEFAULT
'user'` (an integer `role_id` cannot take a subquery default) and lets reads show the
role name without a join.

### Migration `00002_add_roles` (additive — does not edit `00001`)

`00002_add_roles.up.sql`:

```sql
CREATE TABLE IF NOT EXISTS roles (
    name VARCHAR(50) PRIMARY KEY
);

INSERT INTO roles (name) VALUES ('admin'), ('user')
ON CONFLICT (name) DO NOTHING;

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS role VARCHAR(50) NOT NULL DEFAULT 'user'
    REFERENCES roles(name);

-- Give the seeded Alice an admin role so there is an admin to test with.
UPDATE users SET role = 'admin' WHERE email = 'alice@example.com';
```

`00002_add_roles.down.sql`:

```sql
ALTER TABLE users DROP COLUMN IF EXISTS role;
DROP TABLE IF EXISTS roles;
```

Note: the existing migration runner applies pending `.up.sql` files in sorted order and
records them in `migration_logs`, so `00002` applies cleanly on top of an already-migrated
database without touching `00001`.

## Code changes

### models / dtos
- `models.User`: add `Role string` with tag `db:"role"`.
- `dtos.UserDTO`: add `Role string` with tag `json:"role"`.

### service
- Map `Role` through in `GetAllUsers`, `GetUserByID`, `CreateUser`, `UpdateUser`.
- `CreateUser`: if `userDTO.Role == ""`, default to `"user"` before passing to the repo.

### repository (`repositories/user_repository.go`)
- `GetAllUsers` / `GetUserById`: add `role` to the `SELECT` column list.
- `CreateUser`: include `role` in the `INSERT` column list and `RETURNING`.
- `UpdateUser`: **does not** modify `role` (it stays out of the `SET`), but `role` is
  still added to `RETURNING` so the response shows the current value. This prevents a
  `user` from escalating their own role via a self `PUT`. Role is assigned only at
  creation (admin-only `POST`) or by seed.

### New `middleware` package (`middleware/auth.go`)
- `AuthMiddleware(service services.UserService) gin.HandlerFunc`
  - Read `X-User-Id` header; parse to int. Missing/invalid → `401`.
  - Load the user via `service.GetUserByID(id)`. Not found / error → `401`.
  - Set `c.Set("userID", id)` and `c.Set("role", user.Role)`; `c.Next()`.
- `RequireAdmin() gin.HandlerFunc`
  - `403` unless context `role == "admin"`.
- `RequireSelfOrAdmin() gin.HandlerFunc`
  - Allow if `role == "admin"`, or if `:id` param equals context `userID`. Otherwise `403`.

### routes (`routes/user_routes.go`)

Apply `AuthMiddleware` to the `/users` group, then per-route guards:

| Route | Guard |
|-------|-------|
| `GET /users/` (list all) | `RequireAdmin` |
| `GET /users/:id` | `RequireSelfOrAdmin` |
| `POST /users/` | `RequireAdmin` |
| `PUT /users/:id` | `RequireSelfOrAdmin` |
| `DELETE /users/:id` | `RequireAdmin` |

`RegisterUserRoutes` will need access to the `UserService` (or the middleware
constructors) so it can build `AuthMiddleware`. `main.go` wiring passes the service in.

## Error contract

| Condition | Status |
|-----------|--------|
| Missing or non-integer `X-User-Id` | `401` |
| `X-User-Id` references an unknown/deleted user | `401` |
| Authenticated but role/ownership not permitted | `403` |
| Valid and permitted | normal `2xx` |

## Behavioral implications

- `GET /users/` (list all) is **admin-only**. A plain `user` gets `403`, because a
  "self-only" rule cannot apply to a list-everyone request.
- A `user` may `GET`/`PUT` only the record whose `id` matches their `X-User-Id`.
- `POST` and `DELETE` are admin-only.

## Testing plan (manual, via the running server)

With Alice (`id 1`) as admin and a `user`-role account (e.g. John, `id 2`):

1. `GET /users/` with `X-User-Id: 1` → 200 (admin lists all).
2. `GET /users/` with `X-User-Id: 2` → 403 (user cannot list all).
3. `GET /users/2` with `X-User-Id: 2` → 200 (self).
4. `GET /users/1` with `X-User-Id: 2` → 403 (user reading another record).
5. `POST /users/` with `X-User-Id: 2` → 403; with `X-User-Id: 1` → 201.
6. `PUT /users/2` with `X-User-Id: 2` → 200 (self); `PUT /users/1` with `X-User-Id: 2` → 403.
7. `DELETE /users/2` with `X-User-Id: 2` → 403; with `X-User-Id: 1` → 200.
8. Any request with no `X-User-Id` → 401.

## Out of scope (YAGNI)

- Real authentication (passwords, login, JWT signing).
- More than two roles or per-permission granularity.
- Role management endpoints (CRUD on the `roles` table).
