
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
       keycap \
       mcg \
       nordic \
       nutcover \
       nutsandbolts \
       phone \
       pillar_holder \
       pool \
       pottery_wheel \
       simple_stl \
       spiral \
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
