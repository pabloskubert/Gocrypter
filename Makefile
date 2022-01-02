
clean:
	rm -rf dist/
	rm -rf build/

build:
	mkdir build/

	mkdir build/windows
	env GOOS=windows go build --trimpath --ldflags "-s -w" -o build/windows/gocrypter.exe main.go 

	mkdir build/linux
	env GOOS=linux go build --trimpath --ldflags "-s -w" -o build/linux/gocrypter.elf main.go 

	cp -r stub build/windows
	cp -r stub build/linux

	mkdir dist/
	tar --directory=build/windows -zcvf dist/gocrypter_1.0.0_windows_amd64.tar.gz . 
	tar --directory=build/linux   -zcvf dist/gocrypter_1.0.0_linux.tar.gz .


test:
	./gocrypter.elf gocrypter.exe

all: clean build