# Roadmap

1. Split sdf into multiple packages:
	- core: SDF primitives (ie- things with real eval functions)
	- shapes: useful higher-level shapes made from primitives
	- render: functions that render SDFs to verious output formats

2. Cleanup the the public API -- the core package should be clean with a stable API:
	- Make private things that aren't needed externally.
	- Make public things that might be needed by external shape
	  libraries or renderers.

3. Clean up error handling:
	- Get rid of panics. They should really be error returns from the
	  functions.
	- Return errors in those places where we currently might only be
	  printing them.
