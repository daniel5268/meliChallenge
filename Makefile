ci:
	@echo "=== ++++++++++++++++++++ ==="
	@echo "=== Executing unit tests ==="
	@echo "=== ++++++++++++++++++++ ==="

	docker-compose up --remove-orphans --build --exit-code-from ci ci
