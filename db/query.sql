-- name: ListDraftRooms :many
SELECT * FROM draft_rooms
ORDER BY created_at DESC;

-- name: GetDraftRoom :one
SELECT * FROM draft_rooms
WHERE id = $1 LIMIT 1;
