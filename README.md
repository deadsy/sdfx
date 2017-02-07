[![GoDoc](https://godoc.org/github.com/deadsy/libusb?status.svg)](https://godoc.org/github.com/deadsy/sdfx/sdf)

# sdfx

A simple CAD package written in Go (https://golang.org/)

 * Objects are modelled with 2d and 3d signed distance functions (SDFs).
 * Objects are defined with Go code.
 * Objects are rendered to an STL file to be viewed and/or 3d printed.

## How To
 1. See the examples.
 2. Write some Go code to define your own object.
 3. Build and run the Go code.
 4. Preview the STL output in an STL viewer (E.g. http://www.meshlab.net/)
 5. Print the STL file if you like it enough.

## Why?
 * SDFs are cool.
 * As a language Golang > OpenSCAD.
 * OpenSCAD can't do automatic filleting and chamfers.
 * SDFs can do auto filleting and chamfering.
 * Very hackable to try out oddball ideas.
 * Who needs an interactive GUI?

## Gallery

![wheel](docs/gallery/wheel.png "Pottery Wheel Casting Pattern")
![core_box](docs/gallery/core_box.png "Pottery Wheel Core Box")
![cylinder_head](docs/gallery/head.png "Cylinder Head")
![extrude1](docs/gallery/extrude1.png "Twisted Extrusions")
![extrude2](docs/gallery/extrude2.png "Scaled and Twisted Extrusions")
![voronoi](docs/gallery/voronoi.png "2D Points Distance Field")

