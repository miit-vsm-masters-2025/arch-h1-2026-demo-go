# Build stage
FROM golang:1.26 AS builder
WORKDIR /src

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/server ./

# Runtime stage
FROM gcr.io/distroless/base-debian12
WORKDIR /
COPY --from=builder /app/server /server

EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/server"]
