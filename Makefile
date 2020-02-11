.PHONY: build
build:
	(cd ./frontend && npm install && yarn build)
	go build

.PHONY: serve
serve:
	./go-react-webapp-bootstrap serve