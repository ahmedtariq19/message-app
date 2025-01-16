
FROM golang:1.22.5-alpine AS builder
WORKDIR /app

RUN apk --no-cache add curl git

COPY . .

RUN go build -o message-app .


FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/message-app .

EXPOSE 8080

CMD ["./message-app"]