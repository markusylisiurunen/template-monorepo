module github.com/markusylisiurunen/template-monorepo/backend/go-example

go 1.16

require (
	github.com/go-playground/validator/v10 v10.9.0
	github.com/golang-migrate/migrate/v4 v4.14.1
	github.com/google/uuid v1.3.0
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/markusylisiurunen/template-monorepo/package/go/hello v1.0.0
	github.com/spf13/viper v1.9.0
	go.uber.org/zap v1.19.1
)

replace github.com/markusylisiurunen/template-monorepo/package/go/hello => ../../package/go/hello
