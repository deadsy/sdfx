# Features

1. Add 3d bezier surfaces.

2. Fix the distance function for non-linear extrusions. E.g. screws.

3. Generate smaller/better STL files.
See issue #6

4. Add support for 3mf file output.
See: https://github.com/qmuntal/go3mf

5. Add the ability to extrude a 2d SDF along a curve.

6. Add faster evaluation of the SDF for 2d polygons.

7. Add faster evaluation of the SDF for 3d polygons (triangle meshes).
See issue #14.


## General

8. Make the public API of all packages small.
Don't make public symbols that do not need to be.

9. Add more error returns.
All SDF generating functions should return an error, typically used to indicate bad parameters.
These errors should be propagated to the ultimate caller.

10. Panics should be used to indicate fundamental code problems - not just bad parameters.

