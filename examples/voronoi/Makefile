EXEC = $(shell basename $(CURDIR))

all:
	go build

test: all
	./$(EXEC)

clean:
	go clean
	-rm -f *.svg
	-rm -f *.png
	-rm -f *.stl
	-rm -f *.dxf
