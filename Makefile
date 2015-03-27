root := $(shell pwd)
build := $(root)/hdf5
install := $(build)/install

syso := main.syso

all: $(syso)

install: $(syso)
	go install

$(syso): $(install)/lib/libhdf5.a
	mkdir -p $(build)/$@
	cd $(build)/$@ && ar x $(install)/lib/libhdf5.a
	ld -r -o $@ $(build)/$@/*.o

$(install)/lib/libhdf5.a: $(build)/config.log
	$(MAKE) -C $(build) install

$(build)/config.log: $(build)/configure
	cd $(build) && ./configure --prefix=$(install)

$(build)/configure:
	git submodule update --init

clean:
	rm -rf $(syso)
	cd $(build) && (git checkout . && git clean -df)

.PHONY: all install clean
