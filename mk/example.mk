EXEC = $(shell basename $(CURDIR))

.PHONY: all
all:
	go build

.PHONY: test
test: all
	./$(EXEC)
	if [ -f SHA1SUM ]; then shasum -c SHA1SUM; fi;

.PHONY: hash
hash: all
	./$(EXEC)
	$(TOP)/tools/sha1tool.py > SHA1SUM

.PHONY: clean
clean:
	go clean
	-rm -f *.svg
	-rm -f *.png
	-rm -f *.stl
	-rm -f *.dxf
	-rm -f *.3mf
