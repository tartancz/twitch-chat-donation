-- name: CreateDonation :one
insert into donation(user, channel, send_from, amount, text)    
VALUES(?, ?, ?, ?, ?)
RETURNING *;

-- name: GetSumDonationByStreamer :many
SELECT
    CAST(COALESCE(SUM(d.amount), 0) AS INTEGER) AS amount,
    CAST(strftime('%Y-%m-%d', MIN(d."timestamp")) AS TEXT)  AS StartingDate,
    CAST(strftime('%Y-%m-%d', MAX(d."timestamp")) AS TEXT)  AS EndingDate,
    d.channel
FROM donation d 
WHERE d."timestamp" BETWEEN ? AND ?
GROUP BY d.channel;

-- name: GetAllStreamers :many
SELECT DISTINCT d.channel
FROM donation d;

-- name: getLastDonationsByStreamer :many
SELECT d.*
FROM donation d
WHERE d.channel = :channel
ORDER BY d.timestamp DESC
LIMIT :limit OFFSET :offset;

-- name: getLastDonations :many
SELECT d.*
FROM donation d
ORDER BY d.timestamp DESC
LIMIT :limit OFFSET :offset;

