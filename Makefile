VERSION=0.9
DESTDIR?=/
BIN_FOLD=$(DESTDIR)/usr/bin
HOME_FOLD=$(DESTDIR)/var/www/mm

all: gunnify

gunnify: src/github.com/mattn src/github.com/spaghetty/udev src/gunnify.go
	GOPATH=`pwd`; go build src/gunnify.go

src/github.com/mattn:
	GOPATH=`pwd`; go get github.com/mattn/go-gtk/gtk

src/github.com/spaghetty/udev:
	GOPATH=`pwd`; go get github.com/spaghetty/udev
clean:
	echo merda

# rpm: static templates mailproxy mainweb proxymail.conf.example webconf.cfg.example
# 	ln -s . mastermanager-${VERSION}
# 	tar -czpvf mastermanager-${VERSION}.tar.gz \
# 		mastermanager-${VERSION}/src/sip2ser.it \
# 		mastermanager-${VERSION}/static \
# 		mastermanager-${VERSION}/templates \
# 		mastermanager-${VERSION}/proxymail.conf.example \
# 		mastermanager-${VERSION}/webconf.cfg.example \
# 		mastermanager-${VERSION}/getmailrc.example \
# 		mastermanager-${VERSION}/Makefile
# 	rm mastermanager-${VERSION}
# 	mkdir -p ~/rpmbuild/SOURCES/
# 	mv mastermanager-${VERSION}.tar.gz ~/rpmbuild/SOURCES/
# 	rpmbuild -bb mastermanager.spec

# install: mailproxy mainweb
# 	test -d  $(BIN_FOLD) || mkdir -p $(BIN_FOLD)
# 	test -d  $(HOME_FOLD) || mkdir -p $(HOME_FOLD)
# 	install -m 0777 mailproxy $(DESTDIR)/usr/bin
# 	install -m 0777 mainweb $(BIN_FOLD)
# 	install -m 0777 mailproxy $(BIN_FOLD)
# 	cp -R static $(HOME_FOLD)
# 	cp -R templates $(HOME_FOLD)
# 	install -m 0666 proxymail.conf.example $(HOME_FOLD)
# 	install -m 0666 webconf.cfg.example $(HOME_FOLD)
# 	install -m 0666 getmailrc.example $(HOME_FOLD)

