package utils

const (
	SearchUserQuery  = `SELECT username, email, id FROM users WHERE username ILIKE ?`
	SearchGroupQuery = `SELECT * FROM groups WHERE name ILIKE ?`
)
