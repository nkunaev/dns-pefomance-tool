# dns-pefomance-tool

Утилита предназначена для стресс тестирования DNS сервера  
Переменные окружения:

|Key|Default value|Description|
|---|---|---|
|DNS_SERVER|127.0.0.53|DNS server address|
|FQDN_LIST_PATH|./dns_list.txt|Path to file with fqdn to resolve|
|DELAY|2|Requests delay in milliseconds|

# Build

|Command|Description|
|---|---|
|make build|Just build app|
|make all|Run tests before build|

# Dockerize

docker build --no-cache -t dns-stress:$TAG

# Usage

## Local

- Create file `dns_list.txt` with list fqdn to resove
- Choose dns server which will resolve fqdn
  - `export DNS_SERVER=127.0.0.53`
  - `export FQDN_LIST_PATH='./dns_list.txt'`
- Build app or use `go run ./...`

## Kubernetes
- create configmap with fqdn list and mount it to container as file
- set env which described above
- destroy your DNS :)
