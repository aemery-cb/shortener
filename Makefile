server:
	npm run build --prefix=./ui
	go build .

lint:
	golangci-lint run cmd/... pkg/... ui/...