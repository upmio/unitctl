package impl

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/upmio/unitctl/apps/mysql"
	"github.com/upmio/unitctl/apps/sync"
	"github.com/upmio/unitctl/apps/unit"
)

type impl struct {
	db *sql.DB
}

func NewSyncer(username, password, host, database string, port int) (sync.Syncer, error) {
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

func (i *impl) SyncServer(rwHostGroupId, roHostGroupId int, svcGroupName string, mysqlSet unit.MysqlSet) error {
	if i.db == nil {
		return fmt.Errorf("proxysql db is nil")
	}

	_, err := i.db.Exec(cleanHostGroupSql)
	if err != nil {
		return fmt.Errorf("clean hostgroup fail, err: %v", err)
	}

	sqlStr := fmt.Sprintf(insertHostGroupSql, rwHostGroupId, roHostGroupId, svcGroupName)
	_, err = i.db.Exec(sqlStr)
	if err != nil {
		return fmt.Errorf("insert hostgroup fail, err: %v", err)
	}

	_, err = i.db.Exec(cleanServerSql)
	if err != nil {
		return fmt.Errorf("clean server fail, err: %v", err)
	}

	mysql := mysqlSet[0]

	sqlStr = fmt.Sprintf(insertServerSql, rwHostGroupId, mysql.IpAddr, mysql.Port)
	_, err = i.db.Exec(sqlStr)
	if err != nil {
		return fmt.Errorf("insert server fail, err: %v", err)
	}

	_, err = i.db.Exec(loadServerSql)
	if err != nil {
		return fmt.Errorf("load server fail, err: %v", err)
	}

	_, err = i.db.Exec(saveServerSql)
	if err != nil {
		return fmt.Errorf("save server fail, err: %v", err)
	}

	return nil
}

func (i *impl) SyncUsers(defaultHostGroupId, maxConnections int, users mysql.UserSet) error {

	_, err := i.db.Exec(cleanUserSql)
	if err != nil {
		return fmt.Errorf("clean user fail, err: %v", err)
	}

	for _, user := range users {
		sqlStr := fmt.Sprintf(insertUserSql, user.Username, user.Password, defaultHostGroupId, maxConnections)
		_, err := i.db.Exec(sqlStr)
		if err != nil {
			return fmt.Errorf("create user %s fail, err: %v", user.Username, err)
		}
	}

	_, err = i.db.Exec(loadUserSql)
	if err != nil {
		return fmt.Errorf("load user fail, err: %v", err)
	}

	_, err = i.db.Exec(saveUserSql)
	if err != nil {
		return fmt.Errorf("save user fail, err: %v", err)
	}

	return nil
}
