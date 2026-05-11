# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

sdfx is a code-based CAD package written in Go. Objects are modelled as 2D and 3D signed distance functions (SDFs), composed with CSG operations, and rendered to mesh files (STL, 3MF, DXF, SVG, PNG) for viewing or 3D printing. There is no GUI — every model is a Go program (see `examples/*/main.go`) that builds an SDF tree and calls a renderer.

## Common Commands

Top-level Makefile fans out to every `examples/*/` directory:

```
make            # go build every example
make test       # go test ./sdf, then build+run every example and verify its SHA1SUM
make hash       # build+run every example and regenerate SHA1SUM files
make clean      # go clean each example and remove generated stl/svg/png/dxf/3mf
make stldiff    # render all examples at master and at HEAD, diff the STL outputs
                # (use BASE=<ref> HEAD=<ref> to override)
```

Per-example (each `examples/<name>/` includes `mk/example.mk`):

```
cd examples/<name>
make            # go build
make test       # ./<name> then shasum -c SHA1SUM (if present)
make hash       # ./<name> then regenerate SHA1SUM via tools/sha1tool.py
```

Unit tests (separate from the example regression tests):

```
go test ./...                       # run all package tests
go test ./sdf -run TestName         # single test in the core package
```

## Regression Testing Model

There are two layers of regression checking, both driven by rendered output rather than direct geometry assertions:

1. **SHA1SUM files** — each example directory may contain a `SHA1SUM` listing the expected hash of every file the example writes. `make test` runs the example binary and `shasum -c` against it. The committed hashes are synced to amd64 floating-point results, so other architectures may diverge. When a change to core SDF/render code legitimately alters output, regenerate the affected hashes with `make hash` (per example) or top-level `make hash` (all examples).

2. **`make stldiff`** — `tools/stldiff/run.sh` checks out two refs into temp worktrees, builds and runs every example in each, then compares STL outputs and reports IDENTICAL / MINOR (float drift) / MATERIAL (real geometry change) per file. Run this before landing changes to `sdf/` or `render/` to confirm unrelated examples weren't disturbed.

## Architecture

### Core interfaces (`sdf/`)

Everything is built on two interfaces:

```go
type SDF2 interface { Evaluate(p v2.Vec) float64; BoundingBox() Box2 }
type SDF3 interface { Evaluate(p v3.Vec) float64; BoundingBox() Box3 }
```

`Evaluate` returns the signed distance from `p` to the surface (negative = inside). `BoundingBox` is pre-computed at construction time and cached on the struct — every constructor must work it out from its inputs so the renderer can cull space.

The package layers as follows:

- **Primitives** — `Box3D`, `Sphere3D`, `Cylinder3D`, `Circle2D`, `Polygon2D`, bezier/spline profiles, screw threads, text glyphs, gyroid, imported triangle meshes (`mesh3.go`), etc.
- **CSG combinators** — `Union2D/3D`, `Difference2D/3D`, `Intersect2D/3D`, `Cut2D/3D`. Union/Intersect/Difference structs expose a settable min/max blend function so callers can fillet or chamfer joins.
- **Transforms** — `Transform2D/3D` (arbitrary M33/M44), `ScaleUniform`, `RotateUnion`, `RotateCopy`, `Array`, `Offset`, `Shell`.
- **2D → 3D lifts** — `Extrude3D` (with twist/scale variants), `Revolve3D`/`RevolveTheta3D`, `Loft3D`, `Screw3D`.

A model is therefore just a tree of these structs; rendering walks space and calls `Evaluate` on the root, which recurses down the tree.

### Rendering (`render/`)

`Render3` / `Render2` are strategy interfaces (`Render(sdf, output)` + `Info(sdf)`). Top-level helpers `ToSTL`, `To3MF`, `ToDXF`, `ToSVG` wire a renderer to a channel-backed file writer running in its own goroutine.

Available renderers:

- 3D: `NewMarchingCubesUniform(cells)` — fixed grid; `NewMarchingCubesOctree(cells)` — adaptive octree, parallel (`march3p.go`), the current default in examples; `NewMarchingCubesOctreeSingle(cells)` — single-threaded octree; dual-contouring variants live under `render/dc/`.
- 2D: `NewMarchingSquaresUniform(cells)`, `NewMarchingSquaresQuadtree(cells)`, `NewDualContouring2D(cells)`.

The `cells` argument is the resolution along the longest bounding-box axis.

### Higher-level parts (`obj/`)

Parametric, reusable real-world objects (bolts, nuts, knurled heads, gears, geneva drives, gridfinity bins, servo horns, standoffs, panels, etc.) built on top of `sdf`. Each takes a `*Parms` struct and returns `(sdf.SDF3, error)`. Example programs typically combine `sdf` primitives with `obj` parts.

### Vector math (`vec/`)

Small standalone packages: `v2`/`v3` (float64 vectors), `v2i`/`v3i` (integer vectors), `p2` (2D polar), `conv` (conversions between them). Imported throughout as `v2 "github.com/deadsy/sdfx/vec/v2"` etc., so vector literals appear as `v3.Vec{x, y, z}`.

## Conventions (from docs/ROADMAP.md)

- Keep the public API of every package small; do not export symbols that do not need to be.
- Every SDF-generating function should return `(SDFx, error)`. Use the error for bad parameters and propagate it to the ultimate caller.
- Reserve `panic` for fundamental code problems, not parameter validation.
