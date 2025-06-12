.PHONY: run

run-env:
	@echo "Starting Go server using .env file..."
	@set -o allexport; source .env; set +o allexport; go run main.go
	