root := $(shell pwd)
build := $(root)/hdf5
install := $(build)/install

syso := main.syso

all: $(syso)

install: $(syso)
	go install

$(syso): $(install)/lib/libhdf5.a
	cp $< $@

$(install)/lib/libhdf5.a: $(build)/Makefile
	$(MAKE) -C $(build) install

$(build)/Makefile: $(build)/configure
	cd $(build) && ./configure --prefix=$(install)

$(build)/configure:
	git submodule update --init

clean:
	$(MAKE) -C $(build) clean
	rm -rf $(install) $(syso)

.PHONY: all install clean
