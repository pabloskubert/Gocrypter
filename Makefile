
clean:
	rm -rf dist/

build:
	mkdir dist/

	env GOOS=windows go build --trimpath --ldflags "-s -w" -o dist/gocrypter.exe main.go 
	env GOOS=linux go build --trimpath --ldflags "-s -w" -o dist/gocrypter.elf main.go 
	cp -r stub dist/
test:
	./gocrypter.elf gocrypter.exe