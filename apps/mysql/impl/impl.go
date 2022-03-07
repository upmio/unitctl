package impl

import (
	"database/sql"
	"fmt"
	"github.com/upmio/unitctl/apps/mysql"

	_ "github.com/go-sql-driver/mysql"
)

type impl struct {
	db *sql.DB
}

func NewMysqlClient(username, password, host, database string, port int) (mysql.UserClient, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", username, password, host, port, database)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("open %v fail, error: %v", dsn, err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("ping %v fail, error: %v", dsn, err)
	}

	return &impl{
		db: db,
	}, nil
}

func (i *impl) GetMysqlUser(hostIp string) (mysql.UserSet, error) {
	stmt, err := i.db.Prepare(getUserSql)
	if err != nil {
		return nil, fmt.Errorf("prepare stmt %s fail, err: %v", getUserSql, err)
	}

	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			fmt.Printf("close stmt fail, err: %v", err)
		}
	}(stmt)

	rows, err := stmt.Query(hostIp)
	if err != nil {
		return nil, fmt.Errorf("query user fail, err: %v", err)
	}

	var userList = make(mysql.UserSet, 0)

	for rows.Next() {
		user := mysql.NewDefaultUser()
		err := rows.Scan(&user.Username, &user.Username)
		if err != nil {
			return nil, fmt.Errorf("scan user fail, err: %v", err)
		}
		userList = append(userList, user)
	}

	return userList, nil
}
