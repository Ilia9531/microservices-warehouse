package model

// Session — модель для хранения в Redis.
// Хранит связь между SessionUUID и UserUUID.
type Session struct {
	UUID     string `json:"uuid"`
	UserUUID string `json:"user_uuid"`
}
