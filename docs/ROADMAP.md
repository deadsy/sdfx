# Roadmap

1. Split sdf into multiple packages:
	- sdf: SDF primitives (ie- things with real eval functions)
	- shape: useful higher-level shapes made from primitives
	- render: functions that render SDFs to various output formats

2. Cleanup the the public API -- the sdf package should be clean with a stable API:
	- Make private things not needed externally.
	- Make public things needed by external shape libraries or renderers.

3. Clean up error handling:
	- Remove panics as a consequence of invalid function input. They should be error returns from functions.
	- Return errors in places they are currently printed.
