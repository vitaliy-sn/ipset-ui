BIN_DIR             = /usr/local/bin

APP_NAME            = ipset-ui
APP_SERVICE_NAME    = $(APP_NAME).service
APP_BIN_PATH        = $(BIN_DIR)/$(APP_NAME)
DIST_DIR            = web/dist
EMBED_DIST_DIR      = internal/static/dist

VERSION_PKG         := ipset-ui/internal/version
GIT_COMMIT          := $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
ifneq ($(shell git status --porcelain 2>/dev/null),)
GIT_COMMIT          := $(GIT_COMMIT)-dirty
endif
LDFLAGS             := -X $(VERSION_PKG).BuildTime=$(shell date -u +%Y-%m-%dT%H:%M:%SZ) -X $(VERSION_PKG).Commit=$(GIT_COMMIT)

.PHONY: build-frontend build-backend build
.PHONY: dev-backend dev-frontend
.PHONY: install uninstall
.PHONY: docker-compose-up docker-compose-logs

build-frontend:
	@echo "Building frontend..."
	cd web && npm install && npm run build
	@echo "Copying $(DIST_DIR) to $(EMBED_DIST_DIR)..."
	rm -rf $(EMBED_DIST_DIR)
	cp -r $(DIST_DIR) $(EMBED_DIST_DIR)

build-backend: .check-dist
	@echo "Building backend (main server)..."
	go build -ldflags "$(LDFLAGS)" -o $(APP_NAME) ./cmd

.check-dist:
	@echo "Checking embedded frontend ($(EMBED_DIST_DIR))..."
	@test -d $(EMBED_DIST_DIR) || { echo "Error: $(EMBED_DIST_DIR) is missing - run 'make build-frontend' first" >&2; exit 1; }
	@test -n "$$(ls -A $(EMBED_DIST_DIR))" || { echo "Error: $(EMBED_DIST_DIR) is empty - run 'make build-frontend' first" >&2; exit 1; }

build: build-frontend build-backend

dev-backend: build-backend
	LISTEN_ADDRESS=0.0.0.0:6090 \
	FRONTEND_URL=http://localhost:5173 \
	APP_DIR=./test-app-dir \
	sudo -E ./$(APP_NAME)

dev-frontend:
	cd web && npm install && npm run dev

install: build
	sudo install -Dm755 $(APP_NAME) $(APP_BIN_PATH)
	sudo install -Dm644 systemd/$(APP_SERVICE_NAME) /etc/systemd/system/$(APP_SERVICE_NAME)
	sudo systemctl daemon-reload
	sudo systemctl enable $(APP_SERVICE_NAME)
	sudo systemctl restart $(APP_SERVICE_NAME)

uninstall:
	sudo systemctl stop $(APP_SERVICE_NAME) || true
	sudo systemctl disable $(APP_SERVICE_NAME) || true
	sudo rm -f /etc/systemd/system/$(APP_SERVICE_NAME)
	sudo rm -f $(APP_BIN_PATH)
	sudo systemctl daemon-reload

docker-compose-up:
	docker compose up -d --build

docker-compose-logs:
	docker compose logs -f
