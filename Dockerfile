# One image: the Go server serves the API under /api and the built frontend everywhere else
# (docs/adr/0001).

FROM node:24-alpine AS ui
WORKDIR /ui
COPY ui/package.json ui/package-lock.json ./
RUN npm ci
COPY ui/ ./
RUN npm run build

FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=ui /ui/dist ./ui/dist
RUN CGO_ENABLED=0 GOOS=linux go build -o growth-portal ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates \
    && adduser -D -u 10001 app
WORKDIR /app
COPY --from=builder /app/growth-portal .
ENV HOST=0.0.0.0 \
    PORT=8000
USER app
EXPOSE 8000
CMD ["./growth-portal"]
