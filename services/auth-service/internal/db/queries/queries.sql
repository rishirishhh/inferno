-- name: CreateUser :one
INSERT INTO users (
    id,
    email,
    password_hash
)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * 
FROM users
WHERE email = $1;

-- name: UpdateLastLogin :exec
UPDATE users
SET 
    last_login_at = NOW(),
    failed_login_attempts = 0
WHERE id = $1;

-- name: IncrementFailedAttempts :exec
UPDATE users
SET failed_login_attempts = failed_login_attempts +1
WHERE email = $1;

-- name: ChangePassword :exec
UPDATE users
SET
    password_hash = $2, 
    password_changed_at = NOW()
WHERE id = $1;

-- name: CreateRole :one
INSERT INTO roles (
    id,
    name, 
    description
)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetRoleByID :one
SELECT *
FROM roles
WHERE id = $1;

-- name: GetRoleByName :one
SELECT *
FROM roles
WHERE name = $1;

-- name: ListRoles :many
SELECT *
FROM roles
ORDER BY name;


-- name: CreatePermission :one
INSERT INTO permissions(
    id,
    resource,
    action,
    description
)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetPermissionByID :one
SELECT *
FROM permissions
WHERE id = $1;

-- name: GetPermission :one
SELECT *
FROM permissions
WHERE resource = $1
AND action = $2;

-- name: ListPermissions :many
SELECT *
FROM permissions
ORDER BY resource, action;

-- name: AssignRoleToUser :exec
INSERT INTO user_roles (
    user_id,
    role_id
)
VALUES ($1, $2)
ON CONFLICT (user_id, role_id)
DO NOTHING;

-- name: RemoveRoleFromUser :exec
DELETE FROM user_roles
WHERE user_id = $1
AND role_id = $2;

-- name: GetUserRoles :many
SELECT r.*
FROM roles r
JOIN user_roles ur
    ON r.id = ur.role_id
WHERE ur.user_id = $1;

-- name: UserHasRole :one
SELECT EXISTS (
    SELECT 1
    FROM user_roles
    WHERE user_id = $1
    AND role_id = $2
);


-- name: AssignPermissionToRole :exec
INSERT INTO role_permissions (
    role_id,
    permission_id
)
VALUES ($1, $2)
ON CONFLICT (role_id, permission_id)
DO NOTHING;

-- name: RemovePermissionFromRole :exec
DELETE FROM role_permissions
WHERE role_id = $1
AND permission_id = $2;

-- name: GetRolePermissions :many
SELECT p.*
FROM permissions p
JOIN role_permissions rp
    ON p.id = rp.permission_id
WHERE rp.role_id = $1;

-- name: GetUserPermissions :many
SELECT DISTINCT
    p.*
FROM permissions p
JOIN role_permissions rp
    ON p.id = rp.permission_id
JOIN user_roles ur
    ON ur.role_id = rp.role_id
WHERE ur.user_id = $1;


-- name: HasPermission :one
SELECT EXISTS (
    SELECT 1
    FROM permissions p
    JOIN role_permissions rp
        ON p.id = rp.permission_id
    JOIN user_roles ur
        ON ur.role_id = rp.role_id
    WHERE ur.user_id = $1
    AND p.resource = $2
    AND p.action = $3
);