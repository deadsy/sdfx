EXEC = $(shell basename $(CURDIR))

.PHONY: all
all:
	go build

.PHONY: test
test: all
	./$(EXEC)
	if [ -f MD5SUM ]; then md5sum -c MD5SUM; fi;

.PHONY: hash
hash: all
	./$(EXEC)
	$(TOP)/tools/md5tool.py > MD5SUM

.PHONY: clean
clean:
	go clean
	-rm -f *.svg
	-rm -f *.png
	-rm -f *.stl
	-rm -f *.dxf
