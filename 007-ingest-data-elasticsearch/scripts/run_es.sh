docker network create elastic
docker pull docker.elastic.co/elasticsearch/elasticsearch:8.5.0
docker run --name elasticsearch --net elastic -d -p 9200:9200 -p 9300:9300 -e "discovery.type=single-node" -e "xpack.security.enabled=false" -t docker.elastic.co/elasticsearch/elasticsearch:8.5.0