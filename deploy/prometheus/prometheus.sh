#创建conf文件夹
mkdir /home/zzx/GoProject/prometheus/conf
#进入到conf文件夹
cd /home/zzx/GoProject/prometheus/conf
#创建prometheus.yml文件
touch prometheus.yaml
#把如下的内容保存到prometheus.yaml文件中


docker run --name prometheus -d --privileged=true -u=root \
    -p 9090:9090 \
    -v /etc/localtime:/etc/localtime:ro \
    -v /home/zzx/GoProject/prometheus/data:/prometheus/data \
    -v /home/zzx/GoProject/prometheus/conf:/prometheus/config \
    -v /home/zzx/GoProject/prometheus/rules:/prometheus/rules \
    prom/prometheus --config.file=/prometheus/config/prometheus.yaml --web.enable-lifecycle