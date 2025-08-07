.PHONY: run

run:
	echo "Run using infrastructure/"

test: test_integration

test_unit:
	./run-tests.sh

test_integration:
	docker compose -f compose.integration.yml up --build --abort-on-container-exit --exit-code-from test-runner
	docker compose -f compose.integration.yml down
