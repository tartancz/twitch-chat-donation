-- name: CreateDonation :one
insert into donation(user, channel, send_from, amount, text)    
VALUES(?, ?, ?, ?, ?)
RETURNING *;