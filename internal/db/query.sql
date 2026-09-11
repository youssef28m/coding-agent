-- name: CreateChat :one
INSERT INTO chat (name, created_at)
VALUES (?, ?) RETURNING *;

-- name: GetChatByID :one
SELECT id, name, created_at
FROM chat
WHERE id = ?;

-- name: ListChats :many
SELECT id, name, created_at
FROM chat
ORDER BY created_at DESC;

-- name: DeleteChat :exec
DELETE FROM chat
WHERE id = ?;

-- name: CreateMessage :one
INSERT INTO message (chat_id, role, content, created_at)
VALUES (?, ?, ?, ?) RETURNING *;

-- name: GetMessageByID :one
SELECT id, chat_id, role, content, created_at
FROM message
WHERE id = ?;

-- name: ListMessagesByChatID :many
SELECT id, chat_id, role, content, created_at
FROM message
WHERE chat_id = ?
ORDER BY created_at ASC;