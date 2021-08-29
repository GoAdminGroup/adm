GOCMD = go
GOBUILD = $(GOCMD) build
BINARY_NAME = adm
LAST_VERSION = v1.2.23
VERSION = v1.2.24

build:
	GO111MODULE=on $(GOBUILD) -ldflags "-w" -o ./build/mac/$(BINARY_NAME) ./...
	GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o ./build/linux/x86_64/$(BINARY_NAME) ./...
	GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=arm $(GOBUILD) -o ./build/linux/armel/$(BINARY_NAME) ./...
	GO111MODULE=on CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o ./build/windows/x86_64/$(BINARY_NAME).exe ./...
	GO111MODULE=on CGO_ENABLED=0 GOOS=windows GOARCH=386 $(GOBUILD) -o ./build/windows/i386/$(BINARY_NAME).exe ./...
	rm -rf ./build/linux/armel/adm_linux_armel_$(LAST_VERSION).zip
	rm -rf ./build/linux/x86_64/adm_linux_x86_64_$(LAST_VERSION).zip
	rm -rf ./build/windows/x86_64/adm_windows_x86_64_$(LAST_VERSION).zip
	rm -rf ./build/windows/i386/adm_windows_i386_$(LAST_VERSION).zip
	rm -rf ./build/mac/adm_darwin_x86_64_$(LAST_VERSION).zip
	zip -qj ./build/linux/armel/adm_linux_armel_$(VERSION).zip ./build/linux/armel/adm
	zip -qj ./build/linux/x86_64/adm_linux_x86_64_$(VERSION).zip ./build/linux/x86_64/adm
	zip -qj ./build/windows/x86_64/adm_windows_x86_64_$(VERSION).zip ./build/windows/x86_64/adm.exe
	zip -qj ./build/windows/i386/adm_windows_i386_$(VERSION).zip ./build/windows/i386/adm.exe
	zip -qj ./build/mac/adm_darwin_x86_64_$(VERSION).zip ./build/mac/adm
	rm -rf ./build/zip/*
	cp ./build/linux/armel/adm_linux_armel_$(VERSION).zip ./build/zip/
	cp ./build/linux/x86_64/adm_linux_x86_64_$(VERSION).zip ./build/zip/
	cp ./build/windows/x86_64/adm_windows_x86_64_$(VERSION).zip ./build/zip/
	cp ./build/windows/i386/adm_windows_i386_$(VERSION).zip ./build/zip/
	cp ./build/mac/adm_darwin_x86_64_$(VERSION).zip ./build/zip/

.PHONY: build