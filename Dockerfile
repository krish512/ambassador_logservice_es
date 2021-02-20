FROM ubuntu:18.04
ADD ambassador_logservice_es /usr/local/bin/logservice_es
ENTRYPOINT ["/usr/local/bin/logservice_es"]
