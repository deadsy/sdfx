# CCX and CGX

The example results can be consumed by ABAQUS or CalculiX.

CalculiX and CalculiX GraphiX binaries are available for different platforms, like Linux distributions.

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

The above `-i` flag expects a `teapot-hex8.inp` file. The above command creates `frd` files which can be viewed by CalculiX GraphiX:

```bash
cgx teapot-hex8.frd
```
