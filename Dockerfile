FROM golang:alpine AS build
WORKDIR /build
COPY . .
RUN apk add --no-cache make && make build

FROM alpine:3
LABEL author="nkunaev"
WORKDIR /app
USER 65534
COPY ./dns_list.txt dns_list.txt
COPY --from=build /build/out/dns-stress dns-stress
ENTRYPOINT [ "/app/dns-stress" ]
