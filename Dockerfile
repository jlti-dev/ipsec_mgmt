FROM golang:1.16 as builder
WORKDIR /app
COPY go.mod go.mod
RUN go mod download
COPY app/ .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/main .


FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main /app/main
VOLUME /app/config
CMD ["/app/main"]
