package impl

import (
	"database/sql"
	"fmt"
	"github.com/upmio/unitctl/apps/mysql"
	"go.uber.org/zap"

	_ "github.com/go-sql-driver/mysql"
)

type impl struct {
	db     *sql.DB
	logger *zap.SugaredLogger
}

func NewMysqlClient(username, password, host, database string, port int, logger *zap.SugaredLogger) (mysql.UserClient, error) {
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
		db:     db,
		logger: logger,
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
		err := rows.Scan(&user.Username, &user.Password)
		if err != nil {
			return nil, fmt.Errorf("scan user fail, err: %v", err)
		}
		i.logger.Infof("get user %s", user.Username)
		userList = append(userList, user)
	}

	return userList, nil
}
