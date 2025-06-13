local:
	set -a; source .env; set +a && go run .

genmock:
	@mockgen -destination=mock/repository.go -package=mock_repository -source=repository/repository.go

test:
	./coverage.sh

test-report:
	./coverage.sh --html