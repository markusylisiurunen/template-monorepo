module github.com/markusylisiurunen/template-monorepo/backend/go-example

go 1.16

require (
	github.com/go-playground/validator/v10 v10.6.1
	github.com/golang-migrate/migrate/v4 v4.14.1
	github.com/google/uuid v1.2.0
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.7.4
	github.com/markusylisiurunen/template-monorepo/package/go/hello v1.0.0
	github.com/spf13/viper v1.7.1
	go.uber.org/zap v1.16.0
	golang.org/x/tools v0.0.0-20200825202427-b303f430e36d // indirect
)

replace github.com/markusylisiurunen/template-monorepo/package/go/hello => ../../package/go/hello
