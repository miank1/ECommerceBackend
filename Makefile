.PHONY: run-all stop-all logs

run-all:
	@powershell -WindowStyle Hidden -Command "Start-Process powershell -WindowStyle Hidden -ArgumentList '-Command','cd ../user_service; go run cmd/main.go *> service.log'"
	@powershell -WindowStyle Hidden -Command "Start-Process powershell -WindowStyle Hidden -ArgumentList '-Command','cd ../product_service; go run cmd/main.go *> service.log'"
	@powershell -WindowStyle Hidden -Command "Start-Process powershell -WindowStyle Hidden -ArgumentList '-Command','cd ../cart_service; go run cmd/main.go *> service.log'"
	@powershell -WindowStyle Hidden -Command "Start-Process powershell -WindowStyle Hidden -ArgumentList '-Command','cd ../order_service; go run cmd/main.go *> service.log'"
	@powershell -WindowStyle Hidden -Command "Start-Process powershell -WindowStyle Hidden -ArgumentList '-Command','cd ../payment_service; go run cmd/main.go *> service.log'"

	@echo "✅ All services started in background"

stop-all:
	@taskkill /F /IM go.exe > NUL 2>&1
	@echo "🛑 All Go services stopped"

logs:
	@echo "Check logs:"
	@echo "../user_service/service.log"
	@echo "../product_service/service.log"
	@echo "../cart_service/service.log"
	@echo "../order_service/service.log"
	@echo "../payment_service/service.log"