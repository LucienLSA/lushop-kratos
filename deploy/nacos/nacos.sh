docker pull nacos/nacos-server:v2.1.1

mkdir -p /home/zzx/GoProject/data/nacos/{conf,logs,data}

docker run -p 8848:8848 --name nacos -d nacos/nacos-server:v2.1.1

docker cp nacos:/home/nacos/conf /home/zzx/GoProject/data/nacos
docker cp nacos:/home/nacos/data /home/zzx/GoProject/data/nacos
docker cp nacos:/home/nacos/logs /home/zzx/GoProject/data/nacos

chmod 777 /home/zzx/GoProject/data/nacos/{conf,logs,data}

docker rm -f nacos

docker run -d \
-e MODE=standalone \
--privileged=true \
-e JVM_XMS=256m \
-e JVM_XMX=256m \
-e SPRING_DATASOURCE_PLATFORM=mysql \
-e MYSQL_SERVICE_HOST=172.20.187.137 \
-e MYSQL_SERVICE_PORT=3306 \
-e MYSQL_SERVICE_USER=root \
-e MYSQL_SERVICE_PASSWORD=123456 \
-e MYSQL_SERVICE_DB_NAME=nacos \
-e TIME_ZONE='Asia/Shanghai' \
-e NACOS_AUTH_ENABLE=true \
-v /home/zzx/GoProject/data/nacos/logs:/home/nacos/logs \
-v /home/zzx/GoProject/data/nacos/data:/home/nacos/data \
-v /home/zzx/GoProject/data/nacos/conf:/home/nacos/conf \
-p 8848:8848 -p 9848:9848 -p 9849:9849 \
--name nacos --restart=always nacos/nacos-server:v2.1.1