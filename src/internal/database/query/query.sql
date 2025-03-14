--- User Info queries -------------------------------------

-- name: GetAllUsers :many
SELECT * FROM user_info;

-- name: GetUserByID :one
SELECT * FROM user_info WHERE id = $1 LIMIT 1;

-- name: CheckUserExistWithUsernamePassword :one
SELECT *
FROM user_info
WHERE
    username = $1
    AND password = encode (digest ($2, 'sha256'), 'hex')
LIMIT 1;

-- name: GetAllUsersWithPagination :many
SELECT *, COUNT(*) OVER() AS total_count, CEIL(COUNT(*) OVER() / $2::float) AS max_page_id, $2 AS page_size
FROM user_info 
ORDER BY username 
LIMIT $2 OFFSET $1;

-- name: CreateUser :one
INSERT INTO
    user_info (
        email,
        username,
        first_name,
        last_name,
        role,
        password
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        encode(digest($6, 'sha256'), 'hex')
    )
RETURNING
    *;

-- name: UpdateUser :exec
UPDATE user_info
SET
    email = COALESCE($2, email),
    username = COALESCE($3, username),
    first_name = COALESCE($4, first_name),
    last_name = COALESCE($5, last_name),
    password = COALESCE(encode(digest($6, 'sha256'), 'hex'), password)
WHERE
    id = $1;

-- name: DeleteUser :exec
DELETE FROM user_info WHERE id = $1;

--- Posts queries -------------------------------------

-- name: GetAllPosts :many
SELECT * FROM post;

-- name: GetPostByID :one
SELECT * FROM post WHERE id = $1 LIMIT 1;

-- name: GetAllPostsWithPagination :many
SELECT *, COUNT(*) OVER() AS total_count, CEIL(COUNT(*) OVER() / $2::float) AS max_page_id, $2 AS page_size
FROM post
WHERE
    user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;


-- name: CreatePost :one
INSERT INTO post (user_id, title, content) 
VALUES ($1, $2, $3) 
RETURNING *;

-- name: UpdatePost :exec
UPDATE post
SET
    title = COALESCE($2, title),
    content = COALESCE($3, content)
WHERE
    id = $1
    AND user_id = $4;

-- name: DeletePost :exec
DELETE FROM post WHERE id = $1 AND user_id = $2;

-- name: GetTopPostsInPeriodWithPagination :many
SELECT
    p.*,
    prs.total_rating,
    COUNT(*) OVER() AS total_count,
    CEIL(COUNT(*) OVER() / $2::float) AS max_page_id,
    $2 AS page_size
FROM
    post p
JOIN
    post_rating_summary prs ON p.id = prs.post_id
WHERE
    p.created_at >= NOW() - $1::interval
ORDER BY
    prs.total_rating DESC
LIMIT $2 OFFSET $3;

-- name: UpsertRatePost :exec
INSERT INTO
    rating (
        user_id,
        post_id,
        rating_value
    )
VALUES ($1, $2, $3) ON CONFLICT (post_id, user_id) DO
UPDATE
SET
    rating_value = EXCLUDED.rating_value;

-- name: GetAllUserRatedPostsWithPagination :many
SELECT
    p.*,
    r.rating_value,
    COUNT(*) OVER() AS total_count,
    CEIL(COUNT(*) OVER() / $3::float) AS max_page_id,
    $3 AS page_size
FROM
    post p
JOIN
    rating r ON p.id = r.post_id
WHERE
    r.user_id = $1
ORDER BY
    p.created_at DESC
LIMIT $3 OFFSET $2;

-- name: DeleteRatePost :exec
DELETE FROM rating WHERE post_id = $1 AND user_id = $2;

-- Tag queries -------------------------------------

-- name: GetAllTags :many
SELECT * FROM tag;

-- name: GetAllTagsWithPagination :many
SELECT *, COUNT(*) OVER() AS total_count, CEIL(COUNT(*) OVER() / $2::float) AS max_page_id, $2 AS page_size
FROM tag
ORDER BY name
LIMIT $2 OFFSET $1;

-- name: GetTagByID :one
SELECT * FROM tag WHERE id = $1 LIMIT 1;

-- name: CreateTag :one
INSERT INTO tag (name)
VALUES ($1)
RETURNING *;

-- name: UpdateTag :exec
UPDATE tag
SET
    name = COALESCE($2, name)
WHERE
    id = $1;

-- name: DeleteTag :exec
DELETE FROM tag WHERE id = $1;
