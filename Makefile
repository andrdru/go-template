PROJECT_NAME=template
DOCKER_COMPOSE_CMD="docker compose --file docker/docker-compose.yaml"
MIGRATIONFILE = "-- +migrate Up\\n\\n-- +migrate Down"

LOCAL_IMAGE="${PROJECT_NAME}_builder_docker"
LOCAL_IMAGE_ENV="--env EXAMPLE=example"
LOCAL_IMAGE_CMD=docker run --rm \
                   	--volume $(shell pwd):/app \
                   	--workdir /app \
                   	--env HOME="/tmp" \
                   	--network="host" \
                   	"${LOCAL_IMAGE_ENV}" \
                   	"${LOCAL_IMAGE}"

GOLANGCI_LINT_IMAGE="golangci/golangci-lint:v1.53"
GOLANGCI_LINT_CMD=docker run --rm \
                   	--volume $(shell pwd):/app \
                   	--workdir /app \
                   	--env HOME="/tmp" \
                   	"${GOLANGCI_LINT_IMAGE}"

# build local image
#
# make build
# make build FLAGS="--no-cache"
.PHONY: build
build:
	@ sudo docker build --rm ${FLAGS} docker/ -t "${LOCAL_IMAGE}"

 # up start local env
 # sudo required: host.docker.internal "connection refused" issue from container on linux
.PHONY: up
up: down
	@ sudo "${DOCKER_COMPOSE_CMD}" up -d

# down stop local env
.PHONY: down
down:
	@ sudo "${DOCKER_COMPOSE_CMD}" down --remove-orphans

# lint run golangci-lint
.PHONY: lint
lint:
	@ sudo ${GOLANGCI_LINT_CMD} golangci-lint run ./... -v

# gen-go generate all //go:generate
.PHONY: gen-go
gen-go:
	@ go generate -v ./...

# migration create migration file
#
# Usage:
# make migration NAME=example
.PHONY: migration
migration:
	@ if [ -z ${NAME} ]; then echo "Usage: make migration NAME=example"; exit 1; fi
	@ $(eval FILENAME=$$$$(date +%s)"_${NAME}.sql")
	@ echo "$(MIGRATIONFILE)">>migrations/${FILENAME}
	@ git add migrations/${FILENAME}
	@ echo "Created migration file ${FILENAME}"

# db-migrate-up roll up all migrations
.PHONY: db-migrate-up
db-migrate-up: build
	sudo ${LOCAL_IMAGE_CMD} sql-migrate up -config=docker/resources/dbconfig.yml -env=local

# db-migrate-down roll down one last migration
.PHONY: db-migrate-down
db-migrate-down: build
	sudo ${LOCAL_IMAGE_CMD} sql-migrate down -config=docker/resources/dbconfig.yml -env=local

.PHONY: chown
chown:
	@sudo chown ${USER} -R .

# dev-config init dev config example
.PHONY: dev-config
dev-config:
	@ $(eval LOCAL_IMAGE_ENV = "--env IS_DEBUG=true --env HTTP_HOST=0.0.0.0 --env HTTP_PORT=8080 --env POSTGRES_HOST=localhost  --env POSTGRES_PORT=5432 --env POSTGRES_DB=dbname --env POSTGRES_USER=user --env POSTGRES_PASS=pass")
	sudo ${LOCAL_IMAGE_CMD} /bin/sh -c 'envsubst < configs/config.template.yaml > build/config.yaml'
