.PHONY: setup test run clean deploy build

setup:
	docker-compose up -d
	sleep 10
	python scripts/setup_database.py
	python scripts/seed_data.py
	@echo "✅ Development environment ready!"

test:
	pytest tests/unit/ -v
	pytest tests/integration/ -v

test-load:
	docker-compose up -d
	sleep 10
	k6 run tests/load/k6-scenarios.js

run:
	docker-compose up api

clean:
	docker-compose down -v
	docker system prune -f

build:
	docker build -t recommendation-engine/api ./api

deploy:
	docker-compose -f docker-compose.prod.yml up -d

metrics:
	python observability/metrics.py

ab-test:
	python ab-testing/statistical_tests.py

train-models:
	python ml/training/embeddings_generator.py
	python ml/training/collaborative_filtering.py

profile-update:
	python feature-pipeline/user_profile_updater.py