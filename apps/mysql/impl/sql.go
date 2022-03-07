package impl

const (
	getUserSql = `SELECT mysql,authentication_string FROM mysql WHERE host = ? ;`
)
