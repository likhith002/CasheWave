build:
    @go build -o bin/msgbroker .

run: build
    @./bin/msgbroker
