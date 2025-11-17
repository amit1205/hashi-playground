FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o /bin/reporting-worker ./...

FROM alpine:3.20
RUN adduser -D app
USER app
COPY --from=build /bin/reporting-worker /reporting-worker
EXPOSE 8081
ENTRYPOINT ["/reporting-worker"]
