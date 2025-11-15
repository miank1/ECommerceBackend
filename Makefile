SHELL := cmd
SERVICES = userservice productservice orderservice cartservice searchservice paymentservice
BASE_DIR = services
LOG_DIR = logs

GREEN = echo [✔]
RED = echo [✘]
BLUE = echo [➜]

.PHONY: run debug stop clean build seed logs

# -------------------------------
# 🟢 Run all services (with logs)
# -------------------------------
run:
	@$(BLUE) "Starting all Go services..."
	@if not exist $(LOG_DIR) mkdir $(LOG_DIR)
	@for %%s in ($(SERVICES)) do ( \
		$(BLUE) "Starting %%s..." && \
		cd $(BASE_DIR)\%%s\cmd && \
		start "" /B powershell -Command "go run main.go 2>&1 | Tee-Object -Append ..\..\..\$(LOG_DIR)\%%s.log" && \
		cd ..\..\.. \
	)
	@$(GREEN) "All services started."
	@$(BLUE) "Streaming logs (Ctrl+C to exit)..."
	@powershell -Command "Get-Content $(LOG_DIR)\*.log -Wait -Tail 50"

# -------------------------------------
# 🟡 Debug mode: Run ONE service directly
# Usage: make debug service=orderservice
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
# 🟠 Build all services (production)
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
# 🟣 Seed initial database
# ------------------------
seed:
	@$(BLUE) "Running seeders..."
	@cd $(BASE_DIR)\userseeder && go run main.go
	@cd $(BASE_DIR)\productseeder && go run main.go
	@$(GREEN) "Seeding done."

# ------------------------------------
# 🔴 Stop all Go processes cleanly
# ------------------------------------
stop:
	@$(BLUE) "Stopping Go processes..."
	taskkill /FI "IMAGENAME eq go.exe" /T /F >nul 2>&1 || exit 0
	@$(GREEN) "All services stopped."

# -------------------------
# 🧹 Clean logs + builds
# -------------------------
clean:
	@$(BLUE) "Cleaning build and log files..."
	@if exist bin rmdir /s /q bin
	@if exist $(LOG_DIR) rmdir /s /q $(LOG_DIR)
	@$(GREEN) "Clean completed."

# -------------------------
# 📜 Tail all logs manually
# -------------------------
logs:
	@$(BLUE) "Tailing logs (Ctrl+C to stop)..."
	@powershell -Command "Get-Content $(LOG_DIR)\*.log -Wait -Tail 50"
