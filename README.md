# unitctl

## 1.获取configmap对象 并在指定目录生成配置文件
```unitctl get configmap ${configmap} -n ${namespace} -d ${output_dir}```

## 2.获取secret对象，并将data中的内容以json格式输出
```unitctl get secret ${secret} -n ${namespace}```

## 3.根据service group name标签获取本服务组mysql容器ip和port，并且入proxysql本地admin interface中
```unitctl sync server ${svcgroupname} -n ${namespace} --admin-username admin --admin-password 123456 --admin-host 127.0.0.1 --admin-port 6032 --service-type mysql-replication --rw-hostgroup 10 --ro-hostgroup 20```

## 4.根据service group name
```unitctl sync user ${svcgroupname} -n ${namespace} --admin-username admin --admin-password 123456 --admin-host 127.0.0.1 --admin-port 6032 --default-hostgrou 10 --max-connection 10000 --sync-username check --sync-password 123456```