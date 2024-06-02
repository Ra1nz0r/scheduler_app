-- name: GetTask :one
SELECT id,
    date,
    title,
    comment,
    repeat
FROM scheduler
WHERE id = ?
LIMIT 1;
-- name: ListTasks :many
SELECT id,
    date,
    title,
    comment,
    repeat
FROM scheduler
ORDER BY date
LIMIT 10;
-- name: CreateTask :one
INSERT INTO scheduler (date, title, comment, repeat, search)
VALUES (?, ?, ?, ?, ?)
RETURNING *;
-- name: UpdateTask :exec
UPDATE scheduler
set date = ?,
    title = ?,
    comment = ?,
    repeat = ?,
    search = ?
WHERE id = ?;
-- name: UpdateDateTask :exec
UPDATE scheduler
set date = ?
WHERE id = ?;
-- name: DeleteTask :exec
DELETE FROM scheduler
WHERE id = ?;
-- name: SearchTasks :many
SELECT id,
    date,
    title,
    comment,
    repeat
FROM scheduler
WHERE search LIKE ?
ORDER BY date
LIMIT 10;
-- name: SearchDate :many
SELECT id,
    date,
    title,
    comment,
    repeat
FROM scheduler
WHERE date LIKE ?
LIMIT 10;