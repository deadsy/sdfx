# md5tool README

## Overview

We're currently using the sdfx example files as test cases, and we want to ensure that any changes we make to core code don't unintentionally change rendered output.  This tool performs a rudimentary sanity check, using md5 sums to detect changes.  Using md5 sums has the advantage that we only need to store the md5 sums in git, rather than the rendered files.  The total size of all of the `.stl`, `.dxf` etc. files is greater than 1 Gb, so committing all of those to the repo would be bad.

See the next section for what to do if you are getting an error message from md5tool.

## Common Error Messages: 

### `Error: file has changed: <filename>`

The file's content has been changed.  This means that the example code now generates different output than it used to.  You'll need to use `meshlab`, `librecad`, `inkscape`, or your tool of choice to determine whether the new output file is still correct.  If you are sure that the new output is correct, then regenerate the MD5SUM file by running `make update-md5sum` in the example directory.

It's entirely possible and expected that you will occasionally see md5 differences from one version to the next due to the way floating point rounding is handled at various layers of the library and compiler -- we've seen differences in RenderSTL output when changing a variable from `var` to `const`, for instance.  See the Background section below for some other ways to do difference detection for mesh files.

### `Error: file listed in MD5SUM but file not rendered: <filename>`

A file that is listed in the MD5SUM file does not exist in the example directory.  This means that the example code used to render the file, but it didn't do so on the last `make test` run.  If you are sure that the file should no longer be created, then regenerate the MD5SUM file by running `make update-md5sum` in the example directory.

### `Error: file in local filesystem but missing from MD5SUM file: <filename>`

There is a rendered file in the example directory that does not exist in the MD5SUM file. This means that the example code now generates an output file that it didn't used to generate, most likely due to someone adding a new `Render` statement to the example code.  If you are sure that the file should now be rendered, and if you are sure that the file is being rendered correctly, then regenerate the MD5SUM file by running `make update-md5sum` in the example directory.

### `Error: MD5SUM file open/read: <error details>`

There is either no MD5SUM file in the example directory, or we can't open it or make a temporary copy -- see the error details.  In all cases, make sure you have read/write permission to the directory and file, and if the file doesn't exist yet, you can create it by running `make update-md5sum` in the example directory.  You should only need to do the latter if the example is new and has never had an MD5SUM file created yet.

## Caveats and Pitfalls

It's best to run `make clean` before running `make test` to ensure that the md5tool is giving you a full and correct error report.  Otherwise, you are subject to the following minor vulnerability:

- user runs `make test`, creating all of the correct example output files
- user makes a change to code, causing an existing output file to no longer be rendered
- user doesn't run `make clean` -- this means the file still exists 
- user runs `make test` -- this erroneously passes, because the file is still there

We can see several ways that the Makefiles could be altered to detect this condition and run `make clean` automatically, but are holding off on making that change for now out of an abundance of caution.


## Background and Alternatives

The md5 sums appear to be a viable way of detecting changes; as of this writing, sdfx appears to be deterministic, consistently producing the same output for the same model on repeated runs (this is not true of e.g. OpenSCAD).  There is, however, the possibility that a change to sdfx code will cause trivial changes to the rendered output, breaking the md5 sum comparison; we will have to see over time if this creates a problem.

The closest alternative we have found so far that might do a better job than md5 sums would be to use an algorithm that can parse and analyze the rendered output file in a way that is content-aware:  For STL files, for instance, the Hausdorff Distance algorithm (found in meshlab, for instance) might work.  The drawbacks of this approach are that we would need a different analyzer for each rendered file format, and we would either need to keep > 1Gb of reference output files in the repo, or regenerate them from a reference git commit at test runtime.