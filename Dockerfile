FROM golang:alpine AS build
WORKDIR /build
COPY . .
RUN apk add --no-cache make curl && make build

FROM alpine:3
LABEL author="nkunaev"
WORKDIR /app
COPY --from=build /build/out/dns-stress dns-stress
ENTRYPOINT [ "/app/dns-stress" ]