t = 4;               // Thickness
length = 200;        // Length of the beam
width = 25;          // Width of the beam
height = 25 + t + t; // Overall height of the beam

flange_height = t;                       // Height of the top and bottom flanges
web_height = height - 2 * flange_height; // Height of the vertical web

// Top flange
translate([ 0, 0, height - flange_height ]) cube([ length, width, flange_height ]);

// Bottom flange
cube([ length, width, flange_height ]);

// Vertical web
translate([ 0, (width - t) / 2, flange_height ]) cube([ length, t, web_height ]);