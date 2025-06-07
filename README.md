# dns-stress-tool

Утилита предназначена для стресс тестирования DNS сервера  
Переменные окружения:

|Key|Default value|Description|
|---|---|---|
|DNS_SERVER|8.8.8.8|DNS server address|
|FQDN_LIST_PATH|./dns_list.txt|Path to file with fqdn to resolve|
|DELAY|2|Requests delay in milliseconds|
|REQUESTS_AMOUNT|10|Amount iterations and it's size. Example: 10,100,1000,10000 |

# Build

|Command|Description|
|---|---|
|make build|Just build app|
|make all|Run tests before build|
|make dockerize TAG=0.0.1|Build dockerfile with tag. Default TAG='latest'|
|make run|Build container and run it|

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
