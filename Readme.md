
[![CircleCI](https://circleci.com/gh/Lefthander/otus-go-antibruteforce.svg?style=shield)](https://circleci.com/gh/Lefthander/otus-go-antibruteforce)
[![Go Report Card](https://goreportcard.com/badge/github.com/lefthander/otus-go-antibruteforce)](https://goreportcard.com/report/github.com/lefthander/tokenbucket)

# Проектная работа - AntiBruteForce Service

* [Анти-брутфорс описание задания](./anti-bruteforce.md)

## Команды для сборки и тестирования проекта.

* `make generate` - Генерация gRPC API
* `make lint` - Проверка проекта с помощью утилиты golangci-lint
* `make test` - Выполнение Unit тестов
* `make build-server` - Сборка .bin/abf-srv (API Server)
* `make build-client` - Сборка ./bin/abf-ctl (CLI Client for abf-srv)
* `make undeploy` - docker-compose down
* `make deploy`  - docker-compose up -d --build
* `make status`  - docker-compose ps
* `make run` - docker-compose up

## Команды avf-ctl

* `abf-ctl add -n <network> -c=<flag>` - Добавить подсеть в черный - c=false или белый - c=true список
* `abf-ctl del -n <network> -c=<flag>` - Удалить подсеть из черного - c=false или белого - c=true списка
* `abf-ctl show -c=<flag>` - Распечатать весь черный c=false или белый - c=true список
* `abf-ctl test -l <login> -p <password> -i <ip>` - Проверить разрешено ли клиенту с указанными параметрами авторизоваться
* `abf-ctl reset -l <login> -i <ip>` - Сбросить в исходные знаечения TokenBucket'ы для соответсвующего login, password

## Обязательные требования для каждого проекта

* Наличие юнит-тестов на ключевые алгоритмы (ядро) сервиса.
* Докеризация сервиса:
  - наличие валидного Dockerfile для сервиса;
  - сервис и необходимое ему окружение запускаются без ошибок через `make run` (внутри `docker compose up`) в корне проекта,
    при этом все нужные для работы с сервисом порты доступны с хоста.
* Ветка master успешно проходит пайплайн в CI-CD системе (на ваш вкус, Circle CI, Travis CI, Jenkins, GitLab CI и пр.).
**Пайплайн должен в себе содержать**:
  - запуск `go fmt` и проверку, что команда ничего не возвращает (все файлы уже отформатированы);
  - запуск `go vet`;
  - запуск последней версии `golangci-lint` на весь проект с флагом `--enable-all`;
  - запуск `go test` и `go test -race`;
  - сборку бинаря сервиса для версии Go не ниже 1.12. 

При невыполнении хотя бы одного из требований выше - максимальная оценка за проект **4 из 10 баллов**,
несмотря на, например, полностью написанный код сервиса.

Более подробная разбалловка представлена в описании конкретной темы.

### Использование сторонних библиотек

Допускается, но:

- вы должны иметь представление о том, что происходит внутри.
- точно ли подходит данная библиотека для решения вашей задачи?
- не станет ли библиотека узким местом сервиса?
- не полезнее ли написать функционал, которые вы хотите получить от библиотеки, самому?

---

Для упрощения проверки вашего репозитория, рекомендуем использовать значки GitHub
([GitHub badges](https://github.com/dwyl/repo-badges)), а также [Go Report Card](https://goreportcard.com/).
