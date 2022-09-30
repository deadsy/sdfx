# Features

1. Add 3d bezier surfaces.

2. Fix the distance function for non-linear extrusions. E.g. screws.

3. Generate smaller/better STL files.
See issue #6

4. Add the ability to extrude a 2d SDF along a curve.

5. Add faster evaluation of the SDF for 2d polygons.

6. Add faster evaluation of the SDF for 3d polygons (triangle meshes).
See issue #14.


# General

1. Make the public API of all packages small.
Don't make public symbols that do not need to be.

2. Add more error returns.
All SDF generating functions should return an error, typically used to indicate bad parameters.
These errors should be propagated to the ultimate caller.

3. Panics should be used to indicate fundamental code problems - not just bad parameters.

