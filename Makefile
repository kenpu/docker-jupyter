.PHONY: docker launcher start

docker:
	cd ./docker; docker build -t kenpu/jupyter .

launcher:
	mkdir -p ./bin
	go build -o bin/proxy ./proxy/main.go
	cp ./launch.sh ./bin

start: launcher
	cd ./bin; ./launch.sh $(PWD)/data

clean:
	rm -rf ./bin/
