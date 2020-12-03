
DIRS = $(wildcard ./examples/*/.)

all clean:
	for dir in $(DIRS); do \
		$(MAKE) -C $$dir $@ || exit 1; \
	done

test:
	cd sdf; go test; cd ..
	for dir in $(DIRS); do \
		$(MAKE) -C $$dir $@ || exit 1; \
	done
