.PHONY: build
build:
	(cd ./frontend && yarn build)
	go build

.PHONY: serve
serve:
	./go-react-webapp-bootstrap serve