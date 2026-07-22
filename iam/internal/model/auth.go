package model

type Session struct {
	UUID     string
	UserUUID string
}

type Whoami struct {
	UUID  string
	Login string
	Email string
}
