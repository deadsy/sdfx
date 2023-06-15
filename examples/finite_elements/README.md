# CCX and CGX

The example results can be consumed by ABAQUS or CalculiX.

CalculiX and CalculiX GraphiX binaries are available for different platforms, like Linux distributions.

## openSUSE

openSUSE has [CCX](https://software.opensuse.org/package/ccx) package and also [CGX](https://software.opensuse.org/package/cgx) one.

To install CCX and CGX on openSUSE Leap 15.4 you can run as root:

```bash
zypper addrepo https://download.opensuse.org/repositories/science/15.4/science.repo
zypper refresh
zypper install ccx
zypper install cgx
```

# Visualize `inp` file

To visualize the `inp` file by CalculiX GraphiX:

```bash
cgx -c teapot-hex8.inp
```

# Analyze `inp` file

To run the `inp` files by FEA engines like CalculiX:

```bash
ccx -i teapot-hex8
```

The above `-i` flag expects a `teapot-hex8.inp` file.

The above command creates `frd` files containing the results. They can be viewed by CalculiX GraphiX:

```bash
cgx teapot-hex8.frd
```

The boundary conditions and loads used in the calculation will be available together with the results if you run:

```bash
cgx teapot-hex8.frd teapot-hex8.inp
```

## Math solver

The default CCX math solver is `SPOOLES` which is slow. Apparently `PARDISO` is faster and `PaStiX` is the fastest. But it's needed to build the CCX with `PARDISO` or `PaStiX` math libraries.

### PARDISO

#### Linux executable with the Intel Pardiso Solver

You can download [here](https://www.dropbox.com/s/x8axi53l9dk9w4g/ccx_2.19_MT?dl=1) an executable with the Intel Pardiso solver for x86_64 Linux systems. The executable has all the libraries statically linked into it. So it should run by itself without any dependency. Thanks to [these guys](https://www.feacluster.com/calculix.php).
