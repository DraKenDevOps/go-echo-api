FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/server .

FROM alpine:3.22
RUN apk add --no-cache tzdata ca-certificates \
 && adduser -D -u 10001 app \
 && mkdir -p /app/logs /app/uploads \
 && chown -R app:app /app
WORKDIR /app
COPY --from=build /out/server /app/server
USER app
# ENV ENV_MODE=production PORT=8000 TZ=Asia/Bangkok
EXPOSE 8000
# VOLUME ["/app/logs", "/app/uploads"]
# HEALTHCHECK --interval=30s --timeout=3s CMD wget -qO- http://127.0.0.1:${PORT}/health || exit 1
ENTRYPOINT ["/app/server"]
