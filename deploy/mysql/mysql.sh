docker run -d \
  --name mysql \
  -e MYSQL_ROOT_PASSWORD=123456 \
  -e MYSQL_DATABASE=lushop_user \
  -e MYSQL_USER=lucien \
  -e MYSQL_PASSWORD=123456 \
  -p 3306:3306 \
  -v /home/zzx/mysql:/var/lib/mysql \
  mysql:5.7 \
  --character-set-server=utf8mb4 \
  --collation-server=utf8mb4_unicode_ci