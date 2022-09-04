.PHONY: all build clean package

VERSION=1.0.0-1
ARCH=amd64
SHELL := /bin/bash
OUTPUT_DIR=./bin
DEBSRC_DIR=$(OUTPUT_DIR)/deb

GOARCH=$(ARCH) 
GOARM=7
LDFLAGS=-w

all: build

build:
	$(MAKE) $(OUTPUT_DIR)/exporter-$(ARCH)

clean:
	rm -rf $(OUTPUT_DIR)/*

package: $(OUTPUT_DIR)/exporter-$(ARCH)
	$(MAKE) $(OUTPUT_DIR)/shinobi-exporter-$(ARCH).deb

# #################### Binary files ###################

$(OUTPUT_DIR)/exporter-$(ARCH):
	sed -i 's/@dev/$(VERSION)/' ./version.go
	mkdir -p $(OUTPUT_DIR)
	GOARCH=$(GOARCH) GOARM=$(GOARM) go build -o $(OUTPUT_DIR)/exporter-$(ARCH) -v -mod vendor -ldflags $(LDFLAGS) ./cmd/exporter/main.go
	git reset -q --hard

# # ################## Debian package ##################

$(OUTPUT_DIR)/shinobi-exporter-$(ARCH).deb:
	mkdir -p $(DEBSRC_DIR)
	cp -r ./build/package/debian/* $(DEBSRC_DIR)
	mkdir -p $(DEBSRC_DIR)/shinobi-exporter/usr/bin

	cp $(OUTPUT_DIR)/exporter-$(ARCH) $(DEBSRC_DIR)/shinobi-exporter/usr/bin/shinobi-exporter
	cd $(DEBSRC_DIR)/shinobi-exporter; md5deep -l -o f -r usr -r lib > DEBIAN/md5sums

	export DEB_CONTROL_VERSION=$(VERSION) DEB_CONTROL_ARCH=$(ARCH) DEB_CONTROL_SIZE=`du -s $(DEBSRC_DIR) | cut -f1` && \
		envsubst < ./build/package/debian/shinobi-exporter/DEBIAN/control > $(DEBSRC_DIR)/shinobi-exporter/DEBIAN/control

	cd $(DEBSRC_DIR); fakeroot dpkg-deb --build shinobi-exporter
	mv $(DEBSRC_DIR)/shinobi-exporter.deb $(OUTPUT_DIR)/shinobi-exporter-$(ARCH).deb
