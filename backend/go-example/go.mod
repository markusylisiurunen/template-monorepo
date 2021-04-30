module github.com/markusylisiurunen/template-monorepo/backend/go-example

go 1.16

require (
	github.com/markusylisiurunen/template-monorepo/package/go/hello v1.0.0
	go.uber.org/zap v1.16.0
)

replace github.com/markusylisiurunen/template-monorepo/package/go/hello => ../../package/go/hello
