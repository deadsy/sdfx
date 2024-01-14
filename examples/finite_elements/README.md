# Output

## Nodes and elements

Output is generated as `inp` files for ABAQUS or CalculiX. Nodes and elements are saved in separate files, which are then included inside a single, final `inp` file. This final `inp` file can be opened and viewed using [FreeCAD](https://en.wikipedia.org/wiki/FreeCAD).

## Loading

The final `inp` file can apply a distributed *gravity* load to all elements. Additional point loads can be defined by a JSON file.

## Boundary

In the final `inp` file, the boundary is defined by the point restraints inside a JSON file.

## Modification

Currently, a convenient way for making adjustments is provided. Simply modify the load, restraint, and specification JSON files. Examining the unit tests is helpful to figure out how.

# CCX and CGX

The results can be consumed by ABAQUS or CalculiX.

CalculiX and CalculiX GraphiX binaries are available for different platforms, like Linux distributions.

## openSUSE

openSUSE has [CCX](https://software.opensuse.org/package/ccx) package and also [CGX](https://software.opensuse.org/package/cgx) one.

To install CCX and CGX on openSUSE Leap 15.5 you can run as root:

```bash
zypper addrepo https://download.opensuse.org/repositories/science/15.5/science.repo
zypper refresh
zypper install ccx
zypper install cgx
```

# Visualize `inp` file

To visualize the `inp` file by CalculiX GraphiX:

```bash
cgx -c hex8.inp
```

# Analyze `inp` file

To run the `inp` files by FEA engines like CalculiX:

```bash
ccx -i hex8
```

The above `-i` flag expects a `hex8.inp` file.

The above command creates `frd` files containing the results. They can be viewed by CalculiX GraphiX:

```bash
cgx hex8.frd
```

The boundary conditions and loads used in the calculation will be available together with the results if you run:

```bash
cgx hex8.frd hex8.inp
```

## Math solver

The default CCX math solver is `SPOOLES` which is slow. Apparently `PARDISO` is faster and `PaStiX` is the fastest. But it's needed to build the CCX with `PARDISO` or `PaStiX` math libraries.

### PARDISO

#### Linux executable with the Intel Pardiso Solver

You can download [here](https://www.dropbox.com/s/x8axi53l9dk9w4g/ccx_2.19_MT?dl=1) an executable with the Intel Pardiso solver for x86_64 Linux systems. The executable has all the libraries statically linked into it. So it should run by itself without any dependency. Thanks to [these guys](https://www.feacluster.com/calculix.php).
