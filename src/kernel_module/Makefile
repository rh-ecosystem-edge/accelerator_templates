obj-m += ptemplate_char_dev.o

ifndef KVER
KVER=$(shell uname -r)
endif

ifndef KMODVER
	KMODVER=$(shell git describe HEAD 2>/dev/null || git rev-parse --short HEAD)
endif

all:
	make -C /usr/src/kernels/$(KVER)/ M=$(PWD) EXTRA_CFLAGS=-DKMODVER=\\\"$(KMODVER)\\\" modules
clean:
	make -C /usr/src/kernels/$(KVER)/ M=$(PWD) clean
install:
	install -v -m 755 -d /lib/modules/$(KVER)/
	install -v -m 644 simple_char_dev.ko /lib/modules/$(KVER)/simple_char_dev.ko
	depmod -F /lib/modules/$(KVER)/System.map $(KVER)
