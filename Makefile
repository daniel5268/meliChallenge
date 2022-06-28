start:
	@echo "=== +++++++++++++++ ==="
	@echo "=== starting server ==="
	@echo "=== +++++++++++++++ ==="

	docker-compose up --remove-orphans --build start

build:
	@echo "=== ++++++++++++++ ==="
	@echo "=== starting tests ==="
	@echo "=== ++++++++++++++ ==="

	docker-compose build

test:
	@echo "=== ++++++++++++++ ==="
	@echo "=== starting tests ==="
	@echo "=== ++++++++++++++ ==="

	docker-compose up --remove-orphans --build test
