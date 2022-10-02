build: godeps gogen gobuild

clean:
	go clean

gobuild:
	go build -o ./bin

godeps:
	go get

gogen:
	oapi-codegen -config ./api/types.conf.yaml ./api/swagger.yaml
	oapi-codegen -config ./api/api.conf.yaml ./api/swagger.yaml