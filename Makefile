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
	cd $(build)/$@ && ar x $(install)/lib/libhdf5_hl.a
	ld -r -o $@ $(build)/$@/*.o
	cp $< $@

$(install)/lib/libhdf5.a: $(build)/Makefile
	$(MAKE) -C $(build) install

$(build)/Makefile: $(build)/configure
	cd $(build) && ./configure --prefix=$(install)

$(build)/configure:
	git submodule update --init

clean:
	rm -rf $(syso) $(build)/$(syso) $(install)
	$(MAKE) -C $(build) clean

.PHONY: all install clean
