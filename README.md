# Ambassador LogService for Elasticsearch
This is a LogService plugin for Ambassador to send logs to an Elasticsearch index.

## Environment Variables
- `ELASTICSEARCH_ENDPOINTS` can have comma separated values of multiple elasticsearch cluster nodes, default value is `http://localhost:9200`. Custom value example `http://es01.example.com:9200,http://es02.example.com:9200`
- `ELASTICSEARCH_INDEX` can be used to define Index name, default value is ambassador

## Build Steps
- Build go binary on ubuntu server `go build`
- You can use Dockerfile to create docker image `docker build . -t logservice_es --no-cache`
- Use the deployment.yml file to deploy to Kubernetes as `kubectl -f apply deployment.yml`

## License
MIT

## Related Resources
- [Official Docker Respository](https://hub.docker.com/r/krish512/ambassador_logservice_es)
- [Ambassador LogService Plugin](https://www.getambassador.io/docs/latest/topics/running/services/log-service/)
- [Envoy Access Log Default Format](https://www.envoyproxy.io/docs/envoy/v1.10.0/configuration/access_log#default-format-string)
- [Dummy Envoy ALS Sink](https://github.com/dio/metricsink)
- [EnvoyProxy go-control-plane Access Log Service Protobuf](https://github.com/envoyproxy/go-control-plane/blob/main/envoy/service/accesslog/v2/als.pb.go)
- [Envoy Access Log Service Proto](https://github.com/envoyproxy/envoy/blob/main/api/envoy/service/accesslog/v2/als.proto)
- [Envoy GRPC Access Logs fields](https://www.envoyproxy.io/docs/envoy/latest/api-v2/data/accesslog/v2/accesslog.proto#grpc-access-logs)