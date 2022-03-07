package impl

const (
	getUserSql = `SELECT User,authentication_string FROM mysql.user WHERE host = ? ;`
)
