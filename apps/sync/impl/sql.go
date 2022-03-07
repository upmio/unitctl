package impl

const (
	cleanHostGroupSql = `DELETE FROM mysql_replication_hostgroups;`

	insertHostGroupSql = `INSERT INTO mysql_replication_hostgroups(writer_hostgroup,reader_hostgroup,comment) VALUES (%d,%d,'%s');`

	cleanServerSql = `DELETE FROM mysql_servers;`

	insertServerSql = `INSERT INTO mysql_servers(hostgroup_id,hostname,port) VALUES (%d,'%s',%d);`

	loadServerSql = `LOAD MYSQL SERVERS TO RUNTIME;`

	saveServerSql = `SAVE MYSQL SERVERS TO DISK;`

	insertUserSql = `INSERT INTO mysql_users(username,password,default_hostgroup,max_connections) VALUES ('%s','%s',%d,%d);`

	cleanUserSql = `DELETE FROM mysql_users;`

	loadUserSql = `LOAD MYSQL USERS TO RUNTIME`

	saveUserSql = `SAVE MYSQL USERS TO DISK`
)
