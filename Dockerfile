FROM golang:1.20-alpine as builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go get -u ./...

COPY . .

RUN touch ./example/main.go
RUN GOOS=linux GOARCH=amd64 go build -o ./.bin/app ./example/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/.bin/app .

RUN apk add dumb-init
ENTRYPOINT ["/usr/bin/dumb-init", "--"]

CMD ["./app"]