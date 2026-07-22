package model

import "time"

// User — модель для записи в PostgreSQL.
// Соответствует таблице "users".
type User struct {
	UUID                string    `db:"user_uuid"`
	Login               string    `db:"login"`
	PasswordHash        string    `db:"password_hash"`
	Email               string    `db:"email"`
	NotificationMethods []byte    `db:"notification_methods"`
	CreatedAt           time.Time `db:"created_at"` // Техническое поле БД
}
