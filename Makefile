root := $(shell pwd)
source := $(root)/source
target := $(root)/target

ifeq ($(shell uname -s),Darwin)
extension := dylib
else
extension := so
endif

clibrary := libhdf5.$(extension)
glibrary := main.syso

all: $(glibrary)

install: $(glibrary)
	go install

$(glibrary): $(target)/lib/$(clibrary)
	cp $< $@

$(target)/lib/$(clibrary): $(source)/config.log
	$(MAKE) -C $(source) install

$(source)/config.log: $(source)/configure
	cd $(source) && ./configure --prefix=$(target)

$(source)/configure:
	git submodule update --init

clean:
	rm -rf $(target) $(glibrary)
	cd $(source) && (git checkout . && git clean -df)

.PHONY: all install clean
