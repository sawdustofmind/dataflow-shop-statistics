lint:
	golangci-lint run --new-from-rev=master --config ./.golangci.yml

dep:
	go mod tidy -v
	go mod vendor

api-generate:
	oapi-codegen --config openapi_codegen/types.yaml api/swagger.yaml
	oapi-codegen --config openapi_codegen/user_server.yaml api/swagger.yaml