package xhr

type u struct {
	username string
	password string
}

type User interface {
	Username() string
	Password() string
}

func NewUser(username, password string) User {
	return u{username, password}
}

func (t u) Username() string {
	return t.username
}

func (t u) Password() string {
	return t.password
}
