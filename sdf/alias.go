//-----------------------------------------------------------------------------
/*

Type Aliases

Some types originally defined in sdf have been split out to their own packages.
This set of type aliases allows legacy code to continue building while the code
is converted to use the new types.

*/
//-----------------------------------------------------------------------------

package sdf

import (
	"github.com/deadsy/sdfx/vec/p2"
	v2 "github.com/deadsy/sdfx/vec/v2"
	"github.com/deadsy/sdfx/vec/v2i"
	v3 "github.com/deadsy/sdfx/vec/v3"
	"github.com/deadsy/sdfx/vec/v3i"
)

//-----------------------------------------------------------------------------

// V2i is deprecated. Use v2i.Vec instead.
type V2i = v2i.Vec

// V3i is deprecated. Use v3i.Vec instead.
type V3i = v3i.Vec

// V3 is deprecated. Use v3.Vec instead.
type V3 = v3.Vec

// V2 is deprecated. Use v2.Vec instead.
type V2 = v2.Vec

// P2 is deprecated. Use p2.Vec instead.
type P2 = p2.Vec

// V2Set is deprecated. Use v2.VecSet instead.
type V2Set = v2.VecSet

// V3Set is deprecated. Use v3.VecSet instead.
type V3Set = v3.VecSet

//-----------------------------------------------------------------------------
