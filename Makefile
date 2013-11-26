DEPSDIR= $(CURDIR)/.deps
GOPATH= GOPATH=$(DEPSDIR)
GOATLOOPBACK= $(DEPSDIR)/src/github.com/mediocregopher

local: bin deps
	$(GOPATH) go build -o bin/goat goat.go

release: bin rel deps
	@echo -n "What's the name/version number of this release?: "; \
	read version; \
	mkdir bin/goat-$$version; \
	for platform in darwin freebsd linux windows; do \
		for arch in 386 amd64; do \
			echo $$platform $$arch; \
			$(GOPATH) GOOS=$$platform GOARCH=$$arch go build -o bin/goat-$$version/goat_"$$platform"_"$$arch" goat.go; \
		done; \
	done; \
	cd bin; \
	echo "Tar-ing into rel/goat-$$version.tar.gz"; \
	tar cvzf ../rel/goat-$$version.tar.gz goat-$$version

rel:
	mkdir rel

bin:
	mkdir bin

deps: $(GOATLOOPBACK)
	$(GOPATH) go get launchpad.net/goyaml

$(GOATLOOPBACK):
	mkdir -p $(GOATLOOPBACK)/..
	ln -s $(CURDIR) $(GOATLOOPBACK)
