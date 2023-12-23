FROM golang:latest AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -installsuffix netgo -o fibertel-station-exporter .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/fibertel-station-exporter .
EXPOSE 9420

CMD ["/app/fibertel-station-exporter"]
