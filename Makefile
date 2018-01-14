
DIRS = 3dp_nutbolt \
       benchmark \
       bezier \
       bolt_container \
       camshaft \
       challenge \
       cylinder_head \
       dust_collection \
       extrusion \
       fidget \
       finial \
       gears \
       gas_cap \
       geneva \
       nutsandbolts \
       phone \
       pool \
       pottery_wheel \
       simple_stl \
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
