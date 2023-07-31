## Catalyst Go Skeleton

Go skeleton project. Supports REST, GraphQL, and RabbitMQ messaging out of the box.

## TODO

- [x] Supports MySQL
- [x] Supports MySQL migration
- [x] Supports Postgres migration
- [x] Supports Postgres
- [x] Supports Postgres migration
- [x] Messaging handlers for Direct type exchange
- [x] Messaging handlers for Fanout type exchange
- [x] Graceful shutdown messaging
- [x] Usecase examples
- [x] Unit Testing examples
- [x] Integration testing examples
- [x] Rundeck example
- [ ] Protobuf support
- [ ] Caching with Redis in GraphQL requests examples
- [x] Caching with Redis in Repo
- [x] MongoDB support
- [ ] Gzip Compression
- [x] Conventional Commit checker

## **Usage**

**Run the project**

Create a `.env` file contents from `.env.example`.

```shell
go run main.go
```

or using docker-compose

```shell
docker-compose up -d
```

**Run test**

```
go test ./...
```

## Libraries

- [catalystdk](https://github.com/machtwatch/zerolog) for logging
- [chi](https://github.com/go-chi/chi) for API framework
- [gqlgen](https://github.com/99designs/gqlgen) for graphql framework
- [gorm](https://gorm.io) for database
- [viper](https://github.com/spf13/viper) for env variables
- [testify](https://github.com/stretchr/testify) for testing and mocking
- [goose](https://github.com/pressly/goose) for DB migration
- [redis](https://github.com/go-redis/redis) for cache
- [rundeck](https://github.com/lusis/go-rundeck) for job scheduling
- [catalyst-message](https://dev-message.machtwatch.net/swagger/index.html#) for SMS, email and WhatsApp notification
- [resty](https://github.com/go-resty/resty) for HTTP client
- [rabbitmq](https://www.rabbitmq.com/) for messaging
- [ulid](https://github.com/oklog/ulid) for ULID generator
- [jsoniter](https://github.com/json-iterator/go) for JSON encoder/decoder
- [mockery](https://github.com/vektra/mockery) for interface mocking
- [newrelic](https://docs.newrelic.com/docs/apm/agents/go-agent/get-started/introduction-new-relic-go/) for APM

### ENV

It is set using [Viper](https://github.com/spf13/viper), you can see the implementation on `infrastructure/configuration`

It will set config from `config.json` if it's exist,
otherwise it will set from system's ENV

## Logging

We're using Zerolog for logging.

If ENV DEBUG is set to true, then the log is written to stdout in plaintext,
otherwise, the log is written in JSON.

The JSON format makes it easier for log collection system to process the log from the service, makes it easier to
debug for issues in production using log collection system such as Loki.

The plaintext format makes it easier to debug locally

Examples of logging, can be seen at `logger_test.go` or visit the Zerolog's [README](https://github.com/rs/zerolog)

## Migration

We're using Goose for database migration.

To run migrations, use command:

```bash
make migrate DRIVER=<driver> DB_USERNAME=<user> DB_PASSWORD=<password> DB_HOST=<host> DB_NAME=<db-name> SSL_MODE=<ssl-mode>
```

Example

Postgres

```bash
make migrate DRIVER=postgres DB_USERNAME=mamaz DB_PASSWORD=lahachia DB_HOST=127.0.0.1 DB_NAME=example SSL_MODE=disable
```

MySQL

```bash
make migrate DRIVER=mysql DB_USERNAME=mamaz DB_PASSWORD=lahachia DB_HOST=127.0.0.1 DB_NAME=example
```

To create new migration file, use this command:

```bash
make create-migration NAME=<migration_name>
```

This will create a new migration file with content like this

```sql
-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
```

We need to write the SQL ourselves for creating tables and dropping them manually. This is an example:

```sql
-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE "public"."samples" (
    "id" varchar(26) NOT NULL DEFAULT ''::character varying,
    "name" varchar(30) NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz,
    "deleted_at" timestamptz,
    "created_by" varchar(26) DEFAULT NULL::character varying,
    "updated_by" varchar(26) DEFAULT NULL::character varying,
    "deleted_by" varchar(26) DEFAULT NULL::character varying,
    PRIMARY KEY ("id")
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS "public"."samples";
-- +goose StatementEnd

```

Please refer to Goose's [README](https://github.com/pressly/goose) for more documentation.

## Messaging

We're using [RabbitMQ](https://www.rabbitmq.com/) for messaging. It's beeen wrapped to support direct, fanout, and topics exchange types.

Below are several resources to learn RabbitMQ especially the AMQP 0.9.1 protocol:

- https://www.rabbitmq.com/tutorials/amqp-concepts.html
- https://www.cloudamqp.com/blog/part1-rabbitmq-for-beginners-what-is-rabbitmq.html
- https://www.rabbitmq.com/getstarted.html

The core package (infrastructure/messaging/rabitmq/core) has simple tests to just testing the connection to localhost's RabbitMQ server, to run it, you need to execute:

```shell
go test ./... -tags=rabbitmq
```

### How to support graceful retry

To add graceful retry please refer to [this guideline](https://jamtangan.atlassian.net/wiki/spaces/EN/pages/2264727565/How+to+Setup+DLX+with+Parking+Lot+to+Handle+Reliable+Retry).

## Mocking

To easily generate interface mocks for testing, we use [mockery](https://github.com/vektra/mockery).

To generate mock for an interface, use command:

```bash
mockery --dir <interface-directory> --output  <interface-directory>/mocks --name <interface-name> --filename <output-file-name>.go
```

Example:

```bash
mockery --dir ./domain/user --output ./domain/user/mocks --name IUserRepo --filename user.repository_mock.go
```

Usage example:

```go
import (
	"context"
	"testing"
	"github.com/stretchr/testify/assert"
	roleMock "github.com/machtwatch/catalyst-go-skeleton/domain/role/mocks"
	userMock "github.com/machtwatch/catalyst-go-skeleton/domain/user/mocks"
)

func TestUserUC(t *testing.T) {
	userRepoMock := userMock.NewIUserRepo(t)
	roleRepoMock := roleMock.NewIRoleRepo(t)
	userUC := NewUserUC(userRepoMock, roleRepoMock)

	t.Run("it should successfully get user by ID", func(t *testing.T) {
		userRepoMock.On("GetByID", id, context.Background()).Return(userEntity, nil).Once()

		expected := userEntity.MapToUserGQL()
		actual, err := userUC.GetUserByID(id)

		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})
}
```

## Integration Test

We implemented integration test using [testify](https://github.com/stretchr/testify). It's easy to run the test, you can see the implementation on `integration_test`. 

Before running the test, you need to create `integration.test.env` file with the following content:

```shell
# PostgreSQL
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USERNAME=postgres
POSTGRES_PASSWORD=rahasia
POSTGRES_DATABASE=mw_xms_jt_test
POSTGRES_SSL_MODE=disable

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD= # if you have password

# RabbitMQ
RABBITMQ_URI="amqp://guest:guest@localhost:5672/"
RABBITMQ_CONFIG_DURABLE=false
RABBITMQ_CONFIG_EXCLUSIVE=false
RABBITMQ_CONFIG_AUTO_DELETED=false
RABBITMQ_CONFIG_INTERNAL=false
RABBITMQ_CONFIG_NO_WAIT=false

# Hydra
HYDRA_PUBLIC_URL=https://catalyst-hydra-public.machtwatch.net
HYDRA_CLIENT_ID=mw-xms-jt-client
HYDRA_CLIENT_SECRET=01G6MRD8FQ1GKR3W4FG97ZH0BY
HYDRA_AUDIENCE=https://dev-catalyst-auth.machtwatch.net

# JWKS
JWKS_URL=https://catalyst-hydra-public.machtwatch.net/.well-known/jwks.json
JWKS_REFRESH=1
JWKS_TTL=12

# Message Service
NOTIFICATION_SERVICE_BASE_URL=https://dev-message.machtwatch.net
NOTIFICATION_SERVICE_SENDER_EMAIL=no-reply@machtwatch.co.id
NOTIFICATION_SERVICE_SENDER_NAME=No-Reply

# Debug
DEBUG=true
TEST=true
```

After file created, you need to create database and run migrations.

```shell
make migrate DRIVER=postgres DB_USERNAME=postgres DB_PASSWORD=rahasia DB_HOST=localhost DB_NAME=mw_xms_jt_test SSL_MODE=disable
```

Now you can run the test.

```shell
make integration-test
```

## APM (New Relic)

We're using New Relic for Application Performance Monitoring (APM). We implemented newrelic agent using graphql extension so that all queries and mutation will automatically be traced, we also added gorm plugin for database transaction tracing. Usage example:

```
func (repo *UserRepo) GetUserByEmail(ctx context.Context, email string) (userEntity user.UserEmbedded, err error) {
	builder := setSelectUsersEmbedded().Where(sq.Eq{"users.email": email})

	query, args, err := builder.ToSql()
	if err != nil {
		return userEntity, err
	}

	err = repo.db.WithContext(repo.TracingContext(ctx, "SELECT")).Raw(query, args...).Scan(&userEntity).Error
	if err != nil && err == sql.ErrNoRows {
		return userEntity, nil
	}

	if err != nil {
		return userEntity, fmt.Errorf("error on fetching user with username %s, error: %v", email, err)
	}

	return userEntity, nil
}
```

## GraphQL

Install `gqlgen`

```shell
go install github.com/99designs/gqlgen
```

To generate new resolvers from schemas, run the following command

```shell
gqlgen --config config.yml
```

Example

```shell
gqlgen --config graph/v1.yml
```

## Standard Log
The Catalystdk log package can be used directly without needing of any configuration. But we still provide configuration mechanism for flexibility

Configuration

| Field | Type | Description |
|-|-|-|
| AppName | string | application name, needed by log server for easy filtering|
| Level | string | specify log level (default:debug)|
| TimeFormat | string | specify time format (default:RFC3339="2006-01-02T15:04:05Z07:00")|
| Caller | bool | print caller line position or not (default=false)|
| UseJSON| bool | print log in json format (default=false)|

We only need  `SetStdLog` func to configure the log.
The func is not thread safe, only call it when initializing the app.

Initialize Log
```go
import "github.com/macthwatch/catalystdk/go/log"

if err := log.SetStdLog(&log.Config{
		AppName: "catalyst-go-skeleton",
		Caller:  true,
		Level:  "debug", //default is debug
		UseJSON: true,
	}); err != nil {
		// when got error on setting config, it will use the default config.
		log.StdError(context.Background(), nil, err, "init catalystdk log got error")
	}
```

## Rundeck
Before using this rundeck lib, you must create a new job / setup job.
Create the job in rundeck dashboard.
You can contact your team lead / infra team for rundeck dashboard access.
Steps:
- login to rundeck
- create a new project with developer type
- set details and workflow
- in workflow, there is 2 section need setup. options and workflow
	- options : getting value in opts that we send. this variable will we use in workflow
	- workflow : theres many type in workflow in here. in this example we use HTTP Request Node Step
		- remote url : ${option.endpoint}
		- http method : POST
		- headers : {
						"X-Rundeck-Auth": "${option.token}",
						"Content-Type": "application/json"
					}
		- body : {
					"message":"${option.message}"
				 }
- after create it you get job id that will be used in create schedule
- you can check how to create schedule in function CreateSchedule
that's it, or you can check in this documentation for reference https://jamtangan.atlassian.net/wiki/spaces/DEVOPS/pages/1904148481/Rundeck+How+To

## Gitleaks Code Checker
Gitleaks is a tool for detecting and preventing hardcoded secrets like passwords, api keys, and tokens in git repos. The checker runs using the
pre-commit hooks installedby running `make init` when the first time you start develop in this repository.

## Conventional Commit Checker
Conventional commit is a commit message conventions that used to standardize the structure of the team commit messages. The guideline can be
found in [this document](https://jamtangan.atlassian.net/wiki/spaces/EN/pages/2327543902/Conventional+Commit). The checker runs using the
pre-commit hooks installed by running `make init` when the first time you start develop in this repository.