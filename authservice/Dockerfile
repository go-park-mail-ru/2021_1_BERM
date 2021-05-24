FROM golang:alpine as build

COPY . /project

WORKDIR /project

RUN go build -o bin/auth -v ./cmd/

#================================

FROM alpine:latest

COPY --from=build /project/configs /configs

COPY --from=build /project/bin /bin/

EXPOSE 8084 8085

CMD auth
