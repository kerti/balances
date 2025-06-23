local:
	set -a; source .env; set +a && go run .

genmock:
	@mockgen -destination=mock/repository/repository.go -package=mock_repository -source=repository/repository.go
	@mockgen -destination=mock/service/service.go -package=mock_service -source=service/service.go

test:
	./coverage.sh

test-report:
	./coverage.sh --html