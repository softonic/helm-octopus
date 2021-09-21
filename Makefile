.PHONY: dist
dist: export COPYFILE_DISABLE=1 #teach OSX tar to not put ._* files in tar archive
dist: export CGO_ENABLED=0
dist:
	rm -rf build/octopus/* release/*
	mkdir -p build/octopus/bin release/
	cp README.md LICENSE plugin.yaml build/octopus
	GOOS=linux GOARCH=amd64 go build -o build/octopus/bin/octopus -trimpath
	tar -C build/ -zcvf $(CURDIR)/release/helm-octopus-linux-$(VERSION).tgz octopus/
	GOOS=freebsd GOARCH=amd64 go build -o build/octopus/bin/octopus -trimpath
	tar -C build/ -zcvf $(CURDIR)/release/helm-octopus-freebsd-$(VERSION).tgz octopus/
	GOOS=darwin GOARCH=amd64 go build -o build/octopus/bin/octopus -trimpath
	tar -C build/ -zcvf $(CURDIR)/release/helm-octopus-macos-$(VERSION).tgz octopus/
	rm build/octopus/bin/octopus
	GOOS=windows GOARCH=amd64 go build -o build/octopus/bin/octopus.exe -trimpath
	tar -C build/ -zcvf $(CURDIR)/release/helm-octopus-windows-$(VERSION).tgz octopus/