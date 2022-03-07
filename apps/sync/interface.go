package sync

import (
	"github.com/upmio/unitctl/apps/mysql"
	"github.com/upmio/unitctl/apps/unit"
)

type Syncer interface {
	SyncServer(rwHostGroupId, roHostGroupId int, svcGroupName string, mysqlSet unit.MysqlSet) error
	SyncUsers(defaultHostGroupId, maxConnections int, users mysql.UserSet) error
}
