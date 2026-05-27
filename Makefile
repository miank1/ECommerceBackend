SHELL := cmd

SERVICES = userservice productservice orderservice cartservice searchservice paymentservice

BASE_DIR = services
LOG_DIR = logs

NGINX_DIR = C:\tools\nginx-1.31.1
NGINX_CONF = D:/GitHub/ecommerce/nginx/nginx.conf

GREEN = echo [✔]
RED = echo [✘]
BLUE = echo [➜]

.PHONY: run debug stop clean build seed logs start-nginx stop-nginx restart-nginx

# --------------------------------
# 🚀 Start NGINX Gateway
# --------------------------------
start-nginx:
	@$(BLUE) "Starting NGINX..."
	@cd $(NGINX_DIR) && start "" /B powershell -Command ".\nginx.exe -c $(NGINX_CONF)"
	@$(GREEN) "NGINX started."

# --------------------------------
# 🛑 Stop NGINX Gateway
# --------------------------------
stop-nginx:
	@$(BLUE) "Stopping NGINX..."
	@cd $(NGINX_DIR) && powershell -Command ".\nginx.exe -s stop"
	@$(GREEN) "NGINX stopped."

# --------------------------------
# 🔄 Restart NGINX
# --------------------------------
restart-nginx:
	@$(BLUE) "Restarting NGINX..."
	@cd $(NGINX_DIR) && powershell -Command ".\nginx.exe -s reload"
	@$(GREEN) "NGINX reloaded."

# -------------------------------
# 🟢 Run all services + NGINX
# -------------------------------
run: start-nginx
	@$(BLUE) "Starting all Go services..."

	@if not exist $(LOG_DIR) mkdir $(LOG_DIR)

	@for %%s in ($(SERVICES)) do ( \
		$(BLUE) "Starting %%s..." && \
		cd $(BASE_DIR)\%%s\cmd && \
		start "" /B powershell -Command "go run main.go 2>&1 | Tee-Object -Append ..\..\..\$(LOG_DIR)\%%s.log" && \
		cd ..\..\.. \
	)

	@$(GREEN) "All services started."

	@$(BLUE) "Streaming logs..."
	@powershell -Command "Get-Content $(LOG_DIR)\*.log -Wait -Tail 50"

# -------------------------------------
# 🟡 Debug ONE service
# Usage:
# make debug service=orderservice
# -------------------------------------
debug:
	@if "$(service)"=="" ( \
		$(RED) "Usage: make debug service=orderservice" \
	) else ( \
		$(BLUE) "Debugging $(service)..." && \
		cd $(BASE_DIR)\$(service)\cmd && \
		go run main.go \
	)

# ---------------------------------
# 🟠 Build all services
# ---------------------------------
build:
	@$(BLUE) "Building all services..."

	@if not exist bin mkdir bin

	@for %%s in ($(SERVICES)) do ( \
		$(BLUE) "Building %%s..." && \
		cd $(BASE_DIR)\%%s && \
		go build -o ..\..\bin\%%s.exe main.go && \
		cd ..\.. \
	)

	@$(GREEN) "Build completed!"

# ------------------------
# 🌱 Seed databases
# ------------------------
seed:
	@$(BLUE) "Running seeders..."

	@cd $(BASE_DIR)\userseeder && go run main.go
	@cd ..\..

	@cd $(BASE_DIR)\productseeder && go run main.go
	@cd ..\..

	@$(GREEN) "Seeding completed."

# ------------------------------------
# 🔴 Stop all services + NGINX
# ------------------------------------
stop:
	@$(BLUE) "Stopping Go processes..."

	@taskkill /FI "IMAGENAME eq go.exe" /T /F >nul 2>&1 || exit 0

	@$(MAKE) stop-nginx

	@$(GREEN) "All services stopped."

# -------------------------
# 🧹 Clean logs + builds
# -------------------------
clean:
	@$(BLUE) "Cleaning logs and binaries..."

	@if exist bin rmdir /s /q bin
	@if exist $(LOG_DIR) rmdir /s /q $(LOG_DIR)

	@$(GREEN) "Clean completed."

# -------------------------
# 📜 Tail logs
# -------------------------
logs:
	@$(BLUE) "Streaming logs..."
	@powershell -Command "Get-Content $(LOG_DIR)\*.log -Wait -Tail 50"