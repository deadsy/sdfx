
DIRS = 3dp_nutbolt \
       axochord \
       axoloti \
       benchmark \
       bezier \
       bjj \
       bolt_container \
       box \
       camshaft \
       challenge \
       cylinder_head \
       devo \
       dust_collection \
       extrusion \
       fidget \
       finial \
       gears \
       gas_cap \
       geneva \
       nordic \
       nutsandbolts \
       phone \
       pool \
       pottery_wheel \
       simple_stl \
       square_flange \
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
