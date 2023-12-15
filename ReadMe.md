## Распределенное хранилище key:value ##
- Кластеризуется с помощью serf
- Реплицируется через REST API (это нужно доработать и заменить)

Ключи и значения - строковые.

### Запуск кластера в докере ###
    dockebuild . -t distapp
    docker run -p 8080:8080 -e ADVERTISE_ADDR=172.17.0.2 -e CLUSTER_ADDR=172.17.0.2 distapp
    docker run -p 8081:8080 -e ADVERTISE_ADDR=172.17.0.3 -e CLUSTER_ADDR=172.17.0.2 distapp
    docker run -p 8082:8080 -e ADVERTISE_ADDR=172.17.0.4 -e CLUSTER_ADDR=172.17.0.3 distapp