.PHONY: build

build:
	docker build -t mc-docker0 .

run:
	docker run -it mc-docker0