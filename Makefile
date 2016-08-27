.PHONY: docker launcher start

start:
	cd ./bin; ./launch.sh $(PWD)/data

docker:
	cd ./docker; docker build -t kenpu/jupyter .

launcher:
	mkdir -p ./bin
	go build -o bin/proxy ./proxy/main.go
	cp ./launch.sh ./bin

clean:
	rm -rf ./bin/
