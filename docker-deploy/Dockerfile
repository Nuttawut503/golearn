FROM golang:1.17-alpine AS builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -ldflags \
  "-X 'main.buildcommit=$(git log --format="%H" -n 1)' \
  -X 'main.buildtime=$(date "+%A %W %Y %X")'" \ 
  -o app

FROM alpine:latest
WORKDIR /
COPY --from=builder /app/app .
EXPOSE 8080
CMD ["/app"]
