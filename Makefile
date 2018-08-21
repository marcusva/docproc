PREFIX ?= /usr/local
DESTDIR ?=
INSTALL_BIN = install -m 755

VERSION != cat VERSION
GOOS != go env GOOS
GOARCH != go env GOARCH
EXT != if [ "$(GOOS)" = "windows" ]; then echo -n ".exe"; else echo -n ""; fi

APPS := docproc.fileinput docproc.proc docproc.webinput
DISTFILES := LICENSE README.md examples doc/_build/html
DISTNAME := docproc-$(VERSION)-$(GOOS)-$(GOARCH)
DISTDIR := dist/$(DISTNAME)

LDFLAGS := -X main.version=$(VERSION)
STATICS != if [ "$(CGO_ENABLED)" = "0" ]; then echo -n "netgo"; else echo -n ""; fi
TAGS := beanstalk nats nsq $(STATICS)

.PHONY: clean install dist test $(APPS)

all: info $(APPS)

info:
	@echo "Building docproc $(VERSION) for $(GOOS)/$(GOARCH)..."

clean:
	rm -rf dist doc/_build vendor

vendor:
	dep ensure -v

$(DISTDIR):
	mkdir -p $(DISTDIR)

$(DISTDIR)/%:
	go build -tags "$(TAGS)" -ldflags "$(LDFLAGS)" -o $@$(EXT) ./$*

$(APPS): vendor $(DISTDIR) $(APPS:%=$(DISTDIR)/%)

test:
	go test -tags "$(TAGS)" -ldflags "$(LDFLAGS)" ./...

docs:
	make -C doc html

install: $(APPS)
	$(INSTALL_BIN) -d $(DESTDIR)$(PREFIX)/bin
	for app in $(APPS); do \
		$(INSTALL_BIN) $(DISTDIR)/$$app $(DESTDIR)$(PREFIX)/bin; \
	done

dist/$(DISTNAME).tar.gz: docs $(APPS)
	for f in $(DISTFILES); do \
		cp -rf $$f $(DISTDIR); \
	done
	if [ "$(GOOS)" = "windows" ]; then \
		cd dist && zip -r $(DISTNAME).zip $(DISTNAME); \
	else \
		cd dist && tar -czf $(DISTNAME).tar.gz $(DISTNAME); \
	fi

dist: dist/$(DISTNAME).tar.gz

PLATFORMS := windows linux darwin dragonfly freebsd
release: $(PLATFORMS)
$(PLATFORMS): %:
	GOOS=$* make && make dist
