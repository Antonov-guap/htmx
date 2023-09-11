# получение текущего бранча и тега из гит-репозитория
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
GIT_TAG := $(shell git describe --tags --exact-match 2>/dev/null)

# настройка докер-тега для собранного образа
# если текущий коммит помечен тегом, то используем его,
# иначе, если текущий бранч main, то тег latest, иначе branch-name
DOCKER_TAG = $(if $(GIT_TAG),$(GIT_TAG),$(if $(filter main,$(GIT_BRANCH)),latest,$(GIT_BRANCH)))

# сборка докер-образа с приложением
build:
	docker build -t registry.mobbtech.com/hack23-hooked/htmx:$(DOCKER_TAG) .

# сборка и публикация образа в реестре
deploy: build
	docker push registry.mobbtech.com/hack23-hooked/htmx:$(DOCKER_TAG)

# сборка и запуск приложения обычном режиме
# из логирования выкинут консул и шина (много спама)
run: build
	export DOCKER_TAG=$(DOCKER_TAG); \
	docker compose -p hooked up 2>&1 | awk '!/^(iqbus|consul)/'

# сборка и запуск приложения в фоновом режиме
run-daemon: build
	export DOCKER_TAG=$(DOCKER_TAG); \
	docker compose -p hooked up -d

# остановка и удаление  приложения, запущенного ранее в фоновом режиме
stop-clean:
	docker compose -p hooked down
