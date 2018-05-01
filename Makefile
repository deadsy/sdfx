
DIRS = 3dp_nutbolt \
       axoloti \
       benchmark \
       bezier \
       bjj \
       bolt_container \
       box \
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
       text \
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
