
DIRS = $(wildcard ./examples/*/.)

all clean hash:
	for dir in $(DIRS); do \
		$(MAKE) -C $$dir $@ || exit 1; \
	done

test:
	cd sdf; go test; cd ..
	for dir in $(DIRS); do \
		$(MAKE) -C $$dir $@ || exit 1; \
	done

# Cross-branch STL regression check. Renders every example at `master`
# and at HEAD, then compares STL outputs. Use before landing a change
# that touches core rendering or SDF code to confirm no unrelated
# example's geometry was disturbed. Override the base/head refs with
# e.g. `make stldiff BASE=v1.0 HEAD=feature-branch`.
BASE ?= master
HEAD ?= HEAD
stldiff:
	./tools/stldiff/run.sh $(BASE) $(HEAD)

.PHONY: all clean hash test stldiff
