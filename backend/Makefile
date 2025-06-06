local:
	set -a; source .env; set +a && go run .

test:
	chmod u+x ./coverage.sh
	./coverage.sh

test-report:
	chmod u+x ./coverage.sh
	./coverage.sh --html