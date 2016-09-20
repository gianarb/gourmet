.PHONY: build clean

build: clean build/gourmet

build/gourmet:
	go build -o build/gourmet .

clean:
	rm -f build/gourmet
