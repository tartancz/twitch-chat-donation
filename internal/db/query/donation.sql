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