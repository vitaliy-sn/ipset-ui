# ixdx/ipset-ui:latest


# Frontend builder

FROM node:22-bookworm AS frontend-builder
WORKDIR /app

COPY web/package*.json ./
RUN npm install

COPY web/ ./
RUN npm run build


# Backend builder

FROM golang:1.25-bookworm AS backend-builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./
RUN rm -rf internal/static/dist
COPY --from=frontend-builder /app/dist ./internal/static/dist
RUN go build -o ./ipset-ui ./cmd/main.go


# Final image

FROM debian:bookworm
WORKDIR /app

RUN apt update && apt install -y ipset whois && rm -rf /var/lib/apt/lists/* && apt clean all

COPY --from=backend-builder /app/ipset-ui /app/ipset-ui

EXPOSE 8080

CMD ["/app/ipset-ui"]
