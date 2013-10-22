local: bin
	go build -o bin/goat src/goat/main/goat.go

release: bin rel
	@echo -n "What's the name/version number of this release?: "; \
	read version; \
	mkdir bin/goat-$$version; \
	for platform in darwin freebsd linux windows; do \
		for arch in 386 amd64; do \
			echo $$platform $$arch; \
			GOOS=$$platform GOARCH=$$arch go build -o bin/goat-$$version/goat_"$$platform"_"$$arch" src/goat/main/goat.go; \
		done; \
	done; \
	cd bin; \
	echo "Tar-ing into rel/goat-$$version.tar.gz"; \
	tar cvzf ../rel/goat-$$version.tar.gz goat-$$version

rel:
	mkdir rel

bin:
	mkdir bin
