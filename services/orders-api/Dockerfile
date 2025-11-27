FROM golang:1.23-alpine AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o /bin/orders-api ./...

FROM alpine:3.20
RUN adduser -D app
USER app
COPY --from=build /bin/orders-api /orders-api
EXPOSE 8080
ENTRYPOINT ["/orders-api"]
