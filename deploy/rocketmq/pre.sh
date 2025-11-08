mkdir -p /home/lucien/data/rocketmq5.3.3/broker/{logs,store,conf}
mkdir -p /home/lucien/data/rocketmq5.3.3/proxy/{logs,conf}
mkdir -p /home/lucien/data/rocketmq5.3.3/dashboard/conf

touch /home/lucien/data/rocketmq5.3.3/broker/conf/broker.conf
touch /home/lucien/data/rocketmq5.3.3/broker/conf/tools.yml

touch /home/lucien/data/rocketmq5.3.3/proxy/conf/rmq-proxy.json

touch /home/lucien/data/rocketmq5.3.3/dashboard/conf/users.properties

touch /home/lucien/data/rocketmq5.3.3/docker-compose.yaml

chown -R 3000:3000 /home/lucien/data/rocketmq5.3.3/


vim /home/lucien/data/rocketmq5.3.3/broker/conf/broker.conf


brokerClusterName=DefaultCluster
brokerName=broker-a
brokerId=0
deleteWhen=04
fileReservedTime=48
brokerRole=ASYNC_MASTER
flushDiskType=ASYNC_FLUSH
# broker 暴露的IP地址
brokerIP1=192.168.185.128
# 开启认证功能
authenticationEnabled=true
authenticationProvider=org.apache.rocketmq.auth.authentication.provider.DefaultAuthenticationProvider
initAuthenticationUser={"username":"lucien","password":"lucien"}
innerClientAuthenticationCredentials={"accessKey":"lucien","secretKey":"lucien"}
authenticationMetadataProvider=org.apache.rocketmq.auth.authentication.provider.LocalAuthenticationMetadataProvider
# 开启授权功能
authorizationEnabled=true
authorizationProvider=org.apache.rocketmq.auth.authorization.provider.DefaultAuthorizationProvider
authorizationMetadataProvider=org.apache.rocketmq.auth.authorization.provider.LocalAuthorizationMetadataProvider
# 兼容 ACL 1.0 的 plain_acl.yml 文件
migrateAuthFromV1Enabled=true



vim /home/lucien/data/rocketmq5.3.3/broker/conf/tools.yml

accessKey: lucien
secretKey: lucien

vim /home/lucien/data/rocketmq5.3.3/proxy/conf/rmq-proxy.json

{
    "rocketMQClusterName": "DefaultCluster",
    "remotingListenPort": 18680,
    "grpcServerPort": 18681,
    "enableACL": true,
    "authenticationEnabled": true,
    "authenticationProvider": "org.apache.rocketmq.auth.authentication.provider.DefaultAuthenticationProvider",
    "authenticationMetadataProvider": "org.apache.rocketmq.proxy.auth.ProxyAuthenticationMetadataProvider",
    "innerClientAuthenticationCredentials": "{\"accessKey\":\"lucien\", \"secretKey\":\"lucien\"}",
    "enableAclRpcHookForClusterMode": true,
    "authorizationEnabled": true,
    "authorizationProvider": "org.apache.rocketmq.auth.authorization.provider.DefaultAuthorizationProvider",
    "authorizationMetadataProvider": "org.apache.rocketmq.proxy.auth.ProxyAuthorizationMetadataProvider",
    "migrateAuthFromV1Enabled": true
}


vim /home/lucien/data/rocketmq5.3.3/dashboard/conf/users.properties

# 配置 dashboard 登录账号密码
admin=lucien

