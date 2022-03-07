package mysql

type UserClient interface {
	GetMysqlUser(hostIp string) (UserSet, error)
}

type User struct {
	Username string
	Password string
}

func NewDefaultUser() *User {
	return &User{}
}

type UserSet []*User
