FROM node:22-alpine AS web-build
WORKDIR /src/app
COPY app/package.json app/package-lock.json ./
RUN npm ci --no-audit --no-fund
COPY app/ ./
RUN npm run build

FROM golang:1.26.5-alpine AS server-build
WORKDIR /src/server
COPY server/go.mod server/go.sum ./
RUN go mod download
COPY server/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/boredgamz .

FROM alpine:3.21
RUN addgroup -S boredgamz && adduser -S -G boredgamz -u 10001 boredgamz
WORKDIR /app
COPY --from=server-build /out/boredgamz /app/boredgamz
COPY --from=web-build /src/app/dist /app/static
ENV APP_ENV=production PORT=3000 STATIC_DIR=/app/static
EXPOSE 3000
USER boredgamz
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 CMD wget -q -O /dev/null http://127.0.0.1:3000/readyz || exit 1
ENTRYPOINT ["/app/boredgamz"]
