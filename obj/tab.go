//-----------------------------------------------------------------------------
/*

Tabs for Connecting Objects.

*/
//-----------------------------------------------------------------------------

package obj

import (
	"github.com/deadsy/sdfx/sdf"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------

type TabParms struct {
	Size      v3.Vec  // size of tab
	Clearance float64 // clearance between male and female elements
	Angled    bool    // tab at 45 degrees (in x-direction)
}

func Tab(k *TabParms, upper bool) (sdf.SDF3, sdf.SDF3, error) {

	if k.Angled {
		panic("TODO")
	}

	if upper {
		size := v3.Vec{k.Size.X + k.Clearance, k.Size.Y + k.Clearance, k.Size.Z}
		env, err := sdf.Box3D(size, 0)
		if err != nil {
			return nil, nil, err
		}
		env = sdf.Transform3D(env, sdf.Translate3d(v3.Vec{0, 0, 0.5 * k.Size.Z}))
		return nil, env, nil
	}

	// lower
	body, err := sdf.Box3D(k.Size, 0)
	if err != nil {
		return nil, nil, err
	}
	body = sdf.Transform3D(body, sdf.Translate3d(v3.Vec{0, 0, 0.5 * k.Size.Z}))
	return body, nil, nil

}

//-----------------------------------------------------------------------------
