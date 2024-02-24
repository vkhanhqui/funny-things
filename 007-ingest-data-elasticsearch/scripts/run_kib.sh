docker pull docker.elastic.co/kibana/kibana:8.5.0
docker run --name kibana --net elastic -p 5601:5601 docker.elastic.co/kibana/kibana:8.5.0