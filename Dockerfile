# ixdx/ipset-ui:latest

# Frontend builder

FROM node:22-alpine AS frontend-builder
WORKDIR /app

COPY web/package*.json ./
RUN npm install

COPY web/ ./
RUN npm run build

# Backend builder

FROM golang:1.24-alpine AS backend-builder
WORKDIR /app

COPY api/go.mod api/go.sum ./
RUN go mod download

COPY api/ ./
RUN go build -o ./ipset-ui ./cmd/main.go

# Final image

FROM ubuntu:24.04
WORKDIR /app

RUN apt update && apt install -y ipset whois && rm -rf /var/lib/apt/lists/* && apt clean all

COPY --from=backend-builder /app/ipset-ui /app/ipset-ui
COPY --from=frontend-builder /app/dist /app/static

EXPOSE 8080

CMD ["/app/ipset-ui"]
