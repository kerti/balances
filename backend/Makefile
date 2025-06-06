local:
	set -a; source .env; set +a && go run .

test:
	./coverage.sh

test-report:
	./coverage.sh --html