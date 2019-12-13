all: gateway deployer

gateway:
	@echo "build gateway"
	go build -o bin/gateway cmd/gateway/main.go
deployer:
	@echo "build deployer"
	go build -o bin/deployer cmd/deployer//main.go

gateway-image:
	@echo "release gateway image"
	./build/build-image.sh gateway

deployer-image:
	@echo "release deployer image"
	./build/build-image.sh deployer
