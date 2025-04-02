# User Management API

Simple HTTP-server in Go with JWT authorization, PostgreSQL and task system (referal, Twitter, Telegram)

## Run

```
docker-compose up --build
```

then, you can access API on

```
http://localhost:8080
```

## Env variables
When run in docker, it automatically passes DSN variable
```
DB_DSN=postgres://user:password@db:5432/testovoe?sslmode=disable
```
When running locally, specify the variable:
```
DB_DSN=postgres://user:password@localhost:5432/testovoe?sslmode=disable
```
