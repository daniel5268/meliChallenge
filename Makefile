start:
	@echo "=== +++++++++++++++ ==="
	@echo "=== starting server ==="
	@echo "=== +++++++++++++++ ==="

	docker-compose up --remove-orphans --build start

test:
	@echo "=== ++++++++++++++ ==="
	@echo "=== starting tests ==="
	@echo "=== ++++++++++++++ ==="

	docker-compose up --remove-orphans --build test
