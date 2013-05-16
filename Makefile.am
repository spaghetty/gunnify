#all: $(bin_SCRIPTS)

# gunnify: src/github.com/mattn src/github.com/spaghetty/udev src/gunnify.go
# 	GOPATH=`pwd`; go build src/gunnify.go

# src/github.com/mattn:
# 	GOPATH=`pwd`; go get github.com/mattn/go-gtk/gtk

# src/github.com/spaghetty/udev:
# 	GOPATH=`pwd`; go get github.com/spaghetty/udev
# clean:
# 	rm gunnify

GO_DEPS = \
	github.com//mattn/go-gtk/gtk \
	github.com/spaghetty/udev

GO_EXECS = gunnify

CLEANFILES = ${GO_EXECS}

gunnify_GOSRC = src/gunnify.go

$(GO_DEPS):
	export GOPATH=$(realpath $(top_srcdir)) ; $(GOLANG) get $@

$(GO_EXECS): $(GO_DEPS)
	test -z $(realpath $(top_srcdir))/bin || mkdir -p $(realpath $(top_srcdir))/bin
	export GOPATH=$(realpath $(top_srcdir) ) ; $(GOLANG) build -o $(realpath $(top_srcdir))/bin/$@ ./$($@_GOSRC)

install: ${GO_EXECS}
	for prog in $(GO_EXECS); do \
	 install -m 755 $$prog ~/bin/; \
	done

bin_SCRIPTS = $(GO_EXECS)
