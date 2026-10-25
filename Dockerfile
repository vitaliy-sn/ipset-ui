# ixdx/ipset-ui:latest


# Frontend builder

FROM node:24.22-bookworm AS frontend-builder
WORKDIR /app

COPY web/package*.json ./
RUN npm install

COPY web/ ./
RUN npm run build


# Backend builder

FROM golang:1.27-bookworm AS backend-builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./
RUN rm -rf internal/static/dist
COPY --from=frontend-builder /app/dist ./internal/static/dist
ARG VERSION_PKG=ipset-ui/internal/version
RUN set -eux; \
    COMMIT="$(git rev-parse --short HEAD 2>/dev/null || echo unknown)"; \
    if [ -n "$(git status --porcelain 2>/dev/null)" ]; then \
        COMMIT="${COMMIT}-dirty"; \
    fi; \
    LDFLAGS="-X ${VERSION_PKG}.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ) -X ${VERSION_PKG}.Commit=${COMMIT}"; \
    go build -ldflags "${LDFLAGS}" -o ./ipset-ui ./cmd/main.go


# Final image

FROM debian:bookworm
WORKDIR /app

RUN apt update && apt install -y ipset whois && rm -rf /var/lib/apt/lists/* && apt clean all

COPY --from=backend-builder /app/ipset-ui /app/ipset-ui

EXPOSE 8080

CMD ["/app/ipset-ui"]
