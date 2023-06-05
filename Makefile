.PHONY: build
build:
	fyne package -os linux -sourceDir ./chewedfeed -name ChewedFeed

.PHONY: install
install: build
	cp ./automated.toml ~/.config/automated.toml
	rm -rf dist/ && mkdir dist/
	tar -xvf ./ChewedFeed.tar.xz -C dist/
	cd dist/ && make user-install
