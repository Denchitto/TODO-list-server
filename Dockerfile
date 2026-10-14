FROM golang:1.26.1-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . ./
RUN go build -o ./application .

FROM alpine:latest AS application
WORKDIR /app
COPY --from=builder /app/application ./
COPY web ./web/
ENTRYPOINT [ "./application" ]
