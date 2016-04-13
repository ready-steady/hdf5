root := $(shell pwd)
source := $(root)/source
target := $(root)/target

clibrary := libhdf5.a
glibrary := main.syso

all: $(glibrary)

install: $(glibrary)
	go install

$(glibrary): $(target)/lib/$(clibrary)
	$(CC) -shared -lz -Wl,-all_load $< -o $@

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
