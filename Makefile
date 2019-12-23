build:
	docker-compose -f ./docker-compose.yml build --parallel

run: build
	docker-compose -f ./docker-compose.yml up --force-recreate --renew-anon-volumes --abort-on-container-exit accounts accounts-db

run-db:
	docker-compose -f ./docker-compose.yml up --force-recreate --renew-anon-volumes --abort-on-container-exit accounts-db

test: build
	docker-compose -f ./docker-compose.yml up -d --force-recreate --renew-anon-volumes
	docker-compose -f ./docker-compose.yml logs -f accounts-test
	docker-compose -f ./docker-compose.yml down
