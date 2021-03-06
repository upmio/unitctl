# unitctl
unitctl 工具组用于在unit中访问kubernetes api 
在使用configmap作为配置文件将应用配置挂载到容器时，经常发生修改了configmap对象后并没有实际落地到容器中，即没有在容器中生效，需要隔一段时间后，kubelet才会进行同步
为了能在容器中实时的访问configmap和secret ，使用client-go封装了一组get方法

### 1.获取configmap对象 并在指定目录生成配置文件
```unitctl get configmap ${configmap} -n ${namespace} -d ${output_dir}```

### 2.获取secret对象，并将data中的内容以json格式输出
```unitctl get secret ${secret} -n ${namespace}```

### 3.根据service group name标签获取本服务组mysql容器ip和port，并且写入proxysql本地admin interface中
```unitctl sync server ${svcgroupname} -n ${namespace} --admin-username admin --admin-password 123456 --admin-host 127.0.0.1 --admin-port 6032 --service-type mysql-replication --rw-hostgroup 10 --ro-hostgroup 20```

### 4.根据service group name标签获取本服务组mysql容器ip和port，并且连接到指定master mysql节点中 获取所有Host为本proxysql容器ip地址的用户 并写入proxysql本地admin interface中
```unitctl sync user ${svcgroupname} -n ${namespace} --admin-username admin --admin-password 123456 --admin-host 127.0.0.1 --admin-port 6032 --default-hostgroup 10 --max-connection 10000 --sync-username check --sync-password 123456```

### 5.根据pod name 获取指定 pod 以dbscale开头的label键值对
```unitctl get pod ${pod_name} -n ${namespace} --show-lables```