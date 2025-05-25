package model

type User struct {
	ID           uint   `db:"id"`
	Username     string `db:"username"`
	PasswordHash string `db:"password_hash"`
}
