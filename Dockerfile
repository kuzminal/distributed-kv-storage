FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./distapp ./cmd/app/main.go

FROM scratch
COPY --from=builder /app/distapp /usr/bin/distapp
ENTRYPOINT [ "/usr/bin/distapp" ]