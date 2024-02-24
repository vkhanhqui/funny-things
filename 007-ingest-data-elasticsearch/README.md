# ingest-data-elasticsearch
This repo helps you ingest data to ElasticSearch, OpenSearch as well

# ENV
```
python 3.9
pyspark 3.2.1
java 1.8.0_202
elasticsearch 8.5.0
elasticsearch 8.5.0
elasticsearch-spark-30_2.12-8.5.0.jar
```


# Set up
```
# init libs
pip install -r requirements.txt

# Spin elasticsearch
./run_es.sh
python ingest_to_opensearch.py

# Optional
# Spin kibana
./run_kib.sh
```