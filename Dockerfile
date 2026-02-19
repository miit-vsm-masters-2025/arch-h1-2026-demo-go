# Build stage
FROM golang:1.26 AS builder
WORKDIR /src

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/server ./

# Runtime stage
FROM alpine:3.20
RUN addgroup -S app && adduser -S -G app app
WORKDIR /
COPY --from=builder /app/server /server

EXPOSE 8080
USER app
ENTRYPOINT ["/server"]
