#===================#
#== Env Variables ==#
#===================#
DOCKER_COMPOSE_FILE ?= docker-compose.yml

#========================#
#== DATABASE MIGRATION ==#
#========================#

migration-up:
	POSTGRES_USER=postgres POSTGRES_PASSWORD=postgres POSTGRES_DB=card_layer docker compose -f ${DOCKER_COMPOSE_FILE} --profile tools run --rm migrate up

migration-down:
	POSTGRES_USER=postgres POSTGRES_PASSWORD=postgres POSTGRES_DB=card_layer docker compose -f ${DOCKER_COMPOSE_FILE} --profile tools run --rm migrate down 1

migration-clear:
	POSTGRES_USER=postgres POSTGRES_PASSWORD=postgres POSTGRES_DB=card_layer docker compose -f ${DOCKER_COMPOSE_FILE} --profile tools run --rm migrate down

migration-create:
	POSTGRES_USER=postgres POSTGRES_PASSWORD=postgres POSTGRES_DB=card_layer docker compose -f ${DOCKER_COMPOSE_FILE} --profile tools run --rm migrate create -ext sql -dir /migrations $(name)

migration-fix:	
	POSTGRES_USER=postgres POSTGRES_PASSWORD=postgres POSTGRES_DB=card_layer docker compose -f ${DOCKER_COMPOSE_FILE} --profile tools run --rm migrate force ${version}

shell-db:
	docker compose -f ${DOCKER_COMPOSE_FILE} exec db psql -U postgres -d postgres
