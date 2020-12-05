EXEC = $(shell basename $(CURDIR))

all:
	go build

test: all
	./$(EXEC)
	cd $(TOP)/tools/md5tool; make
	$(TOP)/tools/md5tool/md5tool

update-md5sums:
	$(TOP)/tools/md5tool/md5tool -update

clean:
	go clean
	-rm -f *.svg
	-rm -f *.png
	-rm -f *.stl
	-rm -f *.dxf
