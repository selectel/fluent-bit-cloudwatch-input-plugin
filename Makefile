build:
	go build -buildmode=c-shared -o cloudwatch-input.so plugin.go

run:
	fluent-bit -e ./cloudwatch-input.so -c config/fluent-bit.yaml

image:
	docker build -t ghcr.io/selectel/fluent-bit-cloudwatch-input-plugin:development .
