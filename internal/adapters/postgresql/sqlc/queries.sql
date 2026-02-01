-- name: ListAuthors :many
SELECT 
* 
FROM 
products;

-- name: ListProductById :one
SELECT
*
FROM 
products 
WHERE id = $1;