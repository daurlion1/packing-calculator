# Stage 1: copy UI (HTML only)
FROM alpine:latest AS ui-builder
WORKDIR /ui
COPY ui ./dist

# Stage 2: build Go backend
FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY go.mod ./
COPY config ./config
COPY cmd ./cmd
COPY internal ./internal
RUN go build -o server ./cmd/server

# Stage 3: final image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server /app/server
COPY --from=builder /app/config /app/config
COPY --from=ui-builder /ui/dist /app/ui
EXPOSE 8080
CMD ["/app/server"]
