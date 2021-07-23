default: clean install

build: build-laf

install: install-laf

build-laf:
	go build -o "./bin" "./laf"

install-laf:
	go install "./laf"

clean:
	if [ -d "./bin" ]; then rm "./bin"; fi
