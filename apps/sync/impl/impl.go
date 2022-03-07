package impl

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/upmio/unitctl/apps/mysql"
	"github.com/upmio/unitctl/apps/sync"
	"github.com/upmio/unitctl/apps/unit"
	"go.uber.org/zap"
)

type impl struct {
	db     *sql.DB
	logger *zap.SugaredLogger
}

func NewSyncer(username, password, host, database string, port int, logger *zap.SugaredLogger) (sync.Syncer, error) {
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

func (i *impl) SyncServer(rwHostGroupId, roHostGroupId int, svcGroupName string, mysqlSet unit.MysqlSet) error {
	if i.db == nil {
		return fmt.Errorf("proxysql db is nil")
	}

	_, err := i.db.Exec(cleanHostGroupSql)
	if err != nil {
		return fmt.Errorf("clean hostgroup fail, err: %v", err)
	}
	i.logger.Info("clean hostgroup success")

	sqlStr := fmt.Sprintf(insertHostGroupSql, rwHostGroupId, roHostGroupId, svcGroupName)
	_, err = i.db.Exec(sqlStr)
	if err != nil {
		return fmt.Errorf("insert hostgroup fail, err: %v", err)
	}
	i.logger.Info("insert hostgroup success")

	_, err = i.db.Exec(cleanServerSql)
	if err != nil {
		return fmt.Errorf("clean server fail, err: %v", err)
	}
	i.logger.Info("clean server success")

	if len(mysqlSet) == 0 {
		i.logger.Info("there is no master server need to sync")
	} else if len(mysqlSet) == 1 {
		mysql := mysqlSet[0]

		sqlStr = fmt.Sprintf(insertServerSql, rwHostGroupId, mysql.IpAddr, mysql.Port)
		_, err = i.db.Exec(sqlStr)
		if err != nil {
			return fmt.Errorf("insert server fail, err: %v", err)
		}
		i.logger.Infof("insert server %s:%d success", mysql.IpAddr, mysql.Port)
	} else {
		return fmt.Errorf("server need to sync more than one")
	}

	_, err = i.db.Exec(loadServerSql)
	if err != nil {
		return fmt.Errorf("load server fail, err: %v", err)
	}
	i.logger.Info("load server success")

	_, err = i.db.Exec(saveServerSql)
	if err != nil {
		return fmt.Errorf("save server fail, err: %v", err)
	}
	i.logger.Info("save server success")

	return nil
}

func (i *impl) SyncUsers(defaultHostGroupId, maxConnections int, users mysql.UserSet) error {

	_, err := i.db.Exec(cleanUserSql)
	if err != nil {
		return fmt.Errorf("clean user fail, err: %v", err)
	}
	i.logger.Info("clean user success")

	for _, user := range users {
		sqlStr := fmt.Sprintf(insertUserSql, user.Username, user.Password, defaultHostGroupId, maxConnections)
		_, err := i.db.Exec(sqlStr)
		if err != nil {
			return fmt.Errorf("create user %s fail, err: %v", user.Username, err)
		}
		i.logger.Infof("insert user %s success", user.Username)
	}

	_, err = i.db.Exec(loadUserSql)
	if err != nil {
		return fmt.Errorf("load user fail, err: %v", err)
	}
	i.logger.Info("load user success")

	_, err = i.db.Exec(saveUserSql)
	if err != nil {
		return fmt.Errorf("save user fail, err: %v", err)
	}
	i.logger.Info("save user success")

	return nil
}
