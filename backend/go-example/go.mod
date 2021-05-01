module github.com/markusylisiurunen/template-monorepo/backend/go-example

go 1.16

require (
	github.com/go-playground/validator/v10 v10.5.0
	github.com/markusylisiurunen/template-monorepo/package/go/hello v1.0.0
	github.com/spf13/viper v1.7.1
	go.uber.org/zap v1.16.0
)

replace github.com/markusylisiurunen/template-monorepo/package/go/hello => ../../package/go/hello
