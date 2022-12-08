migration-up:
	migrate -source file://migration -database "mysql://root:root@tcp(localhost:3306)/cake-store" up   

migration-down:
	migrate -source file://migration -database "mysql://root:root@tcp(localhost:3306)/cake-store" down

coverage-test:
	go test -v -coverpkg=./... -coverprofile=profile.cov ./...
	go tool cover -func profile.cov

opan-api:
	swagger serve api/v1.yaml --flavor=swagger

docker-compose-up-local:
	docker-compose -f docker-compose-dev.yml up -d

docker-compose-down-local:
	docker-compose -f docker-compose-dev.yml down