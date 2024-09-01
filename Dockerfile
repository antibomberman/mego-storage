FROM golang:1.22.5-alpine AS builder
WORKDIR /app
COPY go.mod .
COPY go.sum .

RUN go mod tidy

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o app cmd/storage/main.go

FROM scratch
COPY --from=builder /app/ .
CMD ["/app"]