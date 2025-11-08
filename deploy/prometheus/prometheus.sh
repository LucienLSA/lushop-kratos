#创建conf文件夹
mkdir /home/lucien/data/prometheus/conf
#进入到conf文件夹
cd /home/lucien/data/prometheus/conf
#创建prometheus.yml文件
touch prometheus.yaml
#把如下的内容保存到prometheus.yaml文件中


docker run --name prometheus -d --privileged=true -u=root \
    -p 9090:9090 \
    -v /etc/localtime:/etc/localtime:ro \
    -v /home/lucien/data/prometheus/data:/prometheus/data \
    -v /home/lucien/data/prometheus/conf:/prometheus/conf \
    -v /home/lucien/data/prometheus/rules:/prometheus/rules \
    prom/prometheus --config.file=/prometheus/conf/prometheus.yaml --web.enable-lifecycle