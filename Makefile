BIN_DIR             = /usr/local/bin

APP_NAME            = ipset-ui
APP_SERVICE_NAME    = $(APP_NAME).service
APP_BIN_PATH        = $(BIN_DIR)/$(APP_NAME)

.PHONY: build-frontend build-backend build
.PHONY: run dev dev-frontend
.PHONY: install uninstall install-service uninstall-service
.PHONY: service-install service-uninstall service-restart service-status service-logs
.PHONY: docker-compose-up docker-compose-logs

build-frontend:
	@echo "Building frontend..."
	cd web && npm install && npm run build
	@echo "Copying dist to internal/static..."
	rm -rf internal/static/dist
	cp -r web/dist internal/static/dist


build-backend:
	@echo "Building backend (main server)..."
	go build -o $(APP_NAME) ./cmd

build: build-frontend build-backend

run: build-frontend
	go run ./cmd

dev: build-backend
	LISTEN_ADDRESS=0.0.0.0:6090 \
	FRONTEND_URL=http://localhost:5173 \
	APP_DIR=./test-app-dir \
	sudo -E ./$(APP_NAME)

dev-frontend:
	cd web && npm install && npm run dev

install: build
	sudo install -Dm755 $(APP_NAME) $(APP_BIN_PATH)

uninstall:
	rm -f $(APP_BIN_PATH)

install-service:
	sudo install -Dm644 systemd/$(APP_SERVICE_NAME) /etc/systemd/system/$(APP_SERVICE_NAME)
	sudo systemctl --user daemon-reload
	sudo systemctl --user enable --now $(APP_SERVICE_NAME)

uninstall-service:
	systemctl --user stop $(APP_SERVICE_NAME) || true
	systemctl --user disable $(APP_SERVICE_NAME) || true
	rm -f /etc/systemd/system/$(APP_SERVICE_NAME)

docker-compose-up:
	docker compose up -d --build

docker-compose-logs:
	docker compose logs -f
