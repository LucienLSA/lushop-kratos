docker run -d \
-p 3000:3000 \
--name=grafana \
--restart=always \
-v /home/zzx/GoProject/grafana-storage:/var/lib/grafana \
grafana/grafana