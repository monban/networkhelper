isoUri = http://tinycorelinux.net/14.x/x86/release/
coreIso = Core-14.0.iso
isoBoot = cdrom/boot
isolinuxcfg = $(isoBoot)/isolinux/isolinux.cfg 
work = `pwd`
networkhelper = ./cpio/bin/networkhelper
stockSquashfs = cdrom/boot/core.gz
newSquashfs = cdrom/boot/nhelper.gz

networkhelper.iso : $(isolinuxcfg) $(networkhelper) $(newSquashfs)
	mkisofs -l -J -r -V TC-custom -no-emul-boot \
	-boot-load-size 4 \
	-boot-info-table -b boot/isolinux/isolinux.bin \
	-c boot/isolinux/boot.cat -o networkhelper.iso cdrom

cdrom :
	mkdir -p cdrom
	7z x $(coreIso) -ocdrom -y

$(coreIso) : 
	wget $(isoUri)$(coreIso)

$(networkhelper) : go.sum go.mod cpio
	GOARCH=386 go build -o /tmp/networkhelper .
	sudo cp /tmp/networkhelper $@

go.sum go.mod : main.go
	go mod tidy
	touch go.sum go.mod

cpio : 
	mkdir -p ./cpio
	zcat $(stockSquashfs) |\
	sudo cpio -i -H newc -d -D ./cpio

$(isolinuxcfg) : isolinux.cfg cdrom
	cp ./isolinux.cfg $@

$(newSquashfs) : $(networkhelper)
	cd ./cpio; sudo find |\
	sudo cpio -o -H newc |\
	gzip -2 >\
	 ../$@

.PHONY: clean
clean :
	-rm -rf ./cdrom
	-sudo rm -rf ./cpio
	-rm ./networkhelper.iso
	-rm -rf squashfs
