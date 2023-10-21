isoUri = http://tinycorelinux.net/14.x/x86/release/
coreIso = Core-14.0.iso
isoBoot = cdrom/boot
isolinuxcfg = $(isoBoot)/isolinux/isolinux.cfg 
work = `pwd`
networkhelper = ./squashfs/bin/networkhelper
stockSquashfs = $(isoBoot)/core.gz
newSquashfs = $(isoBoot)/nhelper.gz

networkhelper.iso : $(isolinuxcfg) $(networkhelper) $(newSquashfs)
	@echo ### Creating ISO ###
	mkisofs -l -J -r -V TC-custom -no-emul-boot \
	-boot-load-size 4 \
	-boot-info-table -b boot/isolinux/isolinux.bin \
	-c boot/isolinux/boot.cat -o networkhelper.iso cdrom

$(isoBoot) : | $(coreIso)
	@echo ### Extracting ISO ###
	mkdir -p cdrom
	7z x $(coreIso) -ocdrom -y

$(coreIso) : 
	wget $(isoUri)$(coreIso)

$(networkhelper) : cmd/networkhelper/main.go go.mod .profile
	@echo ### Building binary ###
	GOARCH=386 go build -o /tmp/networkhelper ./cmd/networkhelper
	mkdir -p squashfs/bin
	mkdir -p squashfs/root
	cp .profile squashfs/root
	cp /tmp/networkhelper $@

cpio : 
	mkdir -p ./cpio
	zcat $(stockSquashfs) |\
	fakeroot cpio -i -H newc -d -D ./cpio

$(isolinuxcfg) : isolinux.cfg | $(isoBoot)
	cp ./isolinux.cfg $@

$(newSquashfs) : $(networkhelper)
	@echo ### Creating new initramfs ###
	cd ./squashfs;\
	find -type f -print0 |\
	fakeroot cpio -H newc -ov0 |\
	gzip -2 >\
	 ../$@

.PHONY: clean
clean :
	-rm -rf ./cdrom
	-rm -rf ./cpio
	-rm ./networkhelper.iso
	-rm -rf squashfs
