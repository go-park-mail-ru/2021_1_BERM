FROM golang:alpine as build

COPY . /project

WORKDIR /project

RUN go build -o bin/image -v ./cmd/

#================================

FROM alpine

WORKDIR /

RUN mkdir image

COPY --from=build /project/configs /configs

COPY --from=build /project/bin /bin/

EXPOSE 8083

CMD image
