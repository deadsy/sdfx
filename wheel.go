//-----------------------------------------------------------------------------
/*


 */
//-----------------------------------------------------------------------------

package main

import (
	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------
// overall build controls

const MM_PER_INCH = 25.4
const SCALE = 1.0 / 0.98 // 2% Al shrinkage
const core_print = true  // add the core print to the wheel
const pie_print = true   // create a 1/n pie segment (n = number of webs)

//-----------------------------------------------------------------------------

// dimension scaling
func dim(x float64) float64 {
	return SCALE * x
}

//-----------------------------------------------------------------------------

// draft angles
var draft_angle = sdf.DtoR(4.0)       // standard overall draft
var core_draft_angle = sdf.DtoR(10.0) // draft angle for the core print

// nominal size values (mm)
var wheel_diameter = dim(MM_PER_INCH * 8.0) // total wheel diameter
var hub_diameter = dim(40.0)                // base diameter of central shaft hub
var hub_height = dim(53.0)                  // height of cental shaft hub
var shaft_diameter = dim(21.0)              // 1" target size - reduced for machining allowance
var shaft_length = dim(45.0)                // length of shaft bore
var wall_height = dim(35.0)                 // height of wheel side walls
var wall_thickness = dim(4.0)               // base thickness of outer wheel walls
var plate_thickness = dim(7.0)              // thickness of wheel top plate
var web_width = dim(4.0)                    // base thickness of reinforcing webs
var web_height = dim(25.0)                  // height of reinforcing webs
var core_height = dim(15.0)                 // height of core print
var number_of_webs = 6                      // number of reinforcing webs

// derived values
var wheel_radius = wheel_diameter / 2
var hub_radius = hub_diameter / 2
var shaft_radius = shaft_diameter / 2
var web_length = wheel_radius - (wall_thickness / 2) - shaft_radius

//-----------------------------------------------------------------------------

// build wheel profile
func wheel_profile() *sdf.PolySDF2 {

	/*

	   """build wheel profile"""
	   draft0 = (hub_h - plate_t) * math.tan(draft_angle)
	   draft1 = (wall_h - plate_t) * math.tan(draft_angle)
	   draft2 = wall_h * math.tan(draft_angle)
	   draft3 = core_h * math.tan(core_draft_angle)
	   if core_print:
	     points = [
	       point((0, 0)),
	       point((0, hub_h + core_h)),
	       point((shaft_r - draft3, hub_h + core_h)),
	       point((shaft_r, hub_h)),
	       point((hub_r, hub_h), 5, 2.0),
	       point((hub_r + draft0, plate_t), 5, 2.0),
	       point((wheel_r - wall_t - draft1, plate_t), 5, 2.0),
	       point((wheel_r - wall_t, wall_h), 5, 1.0),
	       point((wheel_r, wall_h), 5, 1.0),
	       point((wheel_r + draft2, 0)),
	     ]
	   else:
	     points = [
	       point((0, 0)),
	       point((0, hub_h - shaft_l)),
	       point((shaft_r, hub_h - shaft_l)),
	       point((shaft_r, hub_h)),
	       point((hub_r, hub_h), 5, 2.0),
	       point((hub_r + draft0, plate_t), 5, 2.0),
	       point((wheel_r - wall_t - draft1, plate_t), 5, 2.0),
	       point((wheel_r - wall_t, wall_h), 5, 1.0),
	       point((wheel_r, wall_h), 5, 1.0),
	       point((wheel_r + draft2, 0)),
	     ]
	   p = polygon(points, closed=False)
	   p.smooth()
	   return p

	*/

	return nil
}

//-----------------------------------------------------------------------------

// build web profile
func web_profile() *sdf.PolySDF2 {

	/*
	   draft = web_h * math.tan(draft_angle)
	   x0 = (2 * web_w) + draft
	   x1 = web_w + draft
	   x2 = web_w
	   points = [
	     point((-x0, 0)),
	     point((-x1, 0), 3, 1.0),
	     point((-x2, web_h), 3, 1.0),
	     point((x2, web_h), 3, 1.0),
	     point((x1, 0), 3, 1.0),
	     point((x0, 0)),
	   ]
	   p = polygon(points, closed=False)
	   p.smooth()
	   return p
	*/
	return nil
}

//-----------------------------------------------------------------------------

// build core profile
func core_profile() *sdf.PolySDF2 {

	/*
	  draft = core_h * math.tan(core_draft_angle)
	  x0 = (2 * web_w) + draft
	  x1 = web_w + draft
	  x2 = web_w
	  points = [
	    point((0, 0)),
	    point((0, core_h + shaft_l)),
	    point((shaft_r, core_h + shaft_l), 3, 2.0),
	    point((shaft_r, core_h)),
	    point((shaft_r - draft, 0)),
	  ]
	  p = polygon(points, closed=True)
	  p.smooth()
	  return p

	*/
	return nil
}

//-----------------------------------------------------------------------------

func wheel() {
}

//-----------------------------------------------------------------------------
