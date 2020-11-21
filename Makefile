
DIRS = 3dp_nutbolt \
       axochord \
       axoloti \
       benchmark \
       bezier \
       bjj \
       bolt_container \
       offset_box \
       panel_box \
       camshaft \
       cap \
       challenge \
       cylinder_head \
       devo \
       dust_collection \
       extrusion \
       fidget \
       finial \
       flask \
       gears \
       gas_cap \
       geneva \
       keycap \
       maixgo \
       mcg \
       midget \
       nordic \
       nutcover \
       nutsandbolts \
       phone \
       pillar_holder \
       pool \
       pottery_wheel \
       simple_stl \
       spiral \
       sprue \
       square_flange \
       test \
       text \
       voronoi \

all:
	for dir in $(DIRS); do \
		$(MAKE) -C ./examples/$$dir $@; \
	done

test:
	cd sdf; go test; cd ..
	for dir in $(DIRS); do \
		$(MAKE) -C ./examples/$$dir $@; \
	done

clean:
	for dir in $(DIRS); do \
		$(MAKE) -C ./examples/$$dir $@; \
	done
