mkdir -p /home/lucien/data/elasticsearch/config
mkdir -p /home/lucien/data/elasticsearch/data
mkdir -p /home/lucien/data/elasticsearch/plugins
chmod 777 -R /home/lucien/data/elasticsearch
echo "http.host: 0.0.0.0" >> /home/lucien/data/elasticsearch/config/elasticsearch.yml


docker run \
--name elasticsearch \
-p 9200:9200 \
-p 9300:9300 \
-v /home/lucien/data/elasticsearch/data:/usr/share/elasticsearch/data \
-v /home/lucien/data/elasticsearch/plugins:/usr/share/elasticsearch/plugins \
-v /home/lucien/data/elasticsearch/config/elasticsearch.yml:/usr/share/elasticsearch/config/elasticsearch.yml \
-e ES_JAVA_OPTS="-Xms128m -Xmx256m" \
-e "discovery.type=single-node" \
-d elasticsearch:7.10.1