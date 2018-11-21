.PHONY: clean build


clean:
	rm -rf bin

build:
	env GOOS=linux GOARCH=amd64 go build -o bin/noDBExp

test:
	go test .

deploy:
	gcloud alpha functions deploy testiram --trigger-http --entry-point=BrezBaze --runtime=go111 --memory=128MB