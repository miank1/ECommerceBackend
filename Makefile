SHELL := cmd
SERVICES = userservice productservice orderservice cartservice searchservice paymentservice
BASE_DIR = services
LOG_DIR = logs

.PHONY: run stop clean build seed

run:
	@echo 🚀 Starting all Go services (showing real-time logs)...
	@if not exist $(LOG_DIR) mkdir $(LOG_DIR)
	@for %%s in ($(SERVICES)) do ( \
		echo ▶️ Starting %%s... && \
		cd $(BASE_DIR)\%%s\cmd && \
		start "" /B cmd /c "go run main.go >> ..\..\..\$(LOG_DIR)\%%s.log 2>&1" && \
		cd ..\..\.. \
	)
	@echo 🪵 Waiting for services to initialize...
	@powershell -Command "Start-Sleep -Seconds 5"
	@echo ✅ All services are up and running!
	@echo 🪵 Tailing logs (press Ctrl+C to stop)...
	@powershell -Command "Get-Content $(LOG_DIR)\*.log -Wait"

build:
	@echo 🏗️ Building all services...
	@if not exist bin mkdir bin
	@for %%s in ($(SERVICES)) do ( \
		echo 🔨 Building %%s... && \
		cd $(BASE_DIR)\%%s && \
		go build -o ..\..\bin\%%s.exe main.go && \
		cd ..\.. \
	)
	@echo ✅ Build complete!

seed:
	@echo 🌱 Running seeders...
	@cd $(BASE_DIR)\userseeder && go run main.go
	@cd $(BASE_DIR)\productseeder && go run main.go
	@echo ✅ Seeding done!

stop:
	@echo 🛑 Stopping all Go service processes...
	taskkill /FI "IMAGENAME eq go.exe" /T /F >nul 2>&1 || exit 0
	@echo ✅ All Go services stopped.

clean:
	@echo 🧹 Cleaning up build and log files...
	@if exist bin rmdir /s /q bin
	@if exist $(LOG_DIR) rmdir /s /q $(LOG_DIR)
	@echo ✅ Clean complete!
