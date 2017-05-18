
DIRS = benchmark \
       bezier \
       bolt_container \
       camshaft \
       challenge \
       cylinder_head \
       extrusion \
       fidget \
       finial \
       gears \
       geneva \
       pottery_wheel \
       nutsandbolts \
       3dp_nutbolt \
       test \
       voronoi \

all:
	for dir in $(DIRS); do \
		$(MAKE) -C ./examples/$$dir $@; \
	done

format:
	goimports -w .

clean:
	for dir in $(DIRS); do \
		$(MAKE) -C ./examples/$$dir $@; \
	done
