// Parameters
length = 200; // Length of the beam
width = 25;  // Width of the beam
height = 25+4+4; // Overall height of the beam
thickness = 4; // Thickness of the beam

// Derived parameters
flange_height = 4; // Height of the top and bottom flanges
web_height = height - 2 * flange_height; // Height of the vertical web

// I-beam structure
module ibeam() {
    // Top flange
    translate([0, 0, height - flange_height])
        cube([length, width, flange_height]);

    // Bottom flange
    cube([length, width, flange_height]);

    // Vertical web
    translate([0, (width - thickness) / 2, flange_height])
        cube([length, thickness, web_height]);
}

// Generate the I-beam
ibeam();