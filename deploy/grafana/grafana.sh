mkdir /home/lucien/data/grafana-storage

chmod 777 -R /home/lucien/data/grafana-storage

docker run -d \
-p 3000:3000 \
--name=grafana \
--restart=always \
-v /home/lucien/data/grafana-storage:/var/lib/grafana \
grafana/grafana