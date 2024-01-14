l = 200;
d = 27.87;
t = 4;

translate([ 0, 0, d / 2 ]) rotate([ 0, 90, 0 ]) difference()
{
    // Outer cylinder
    cylinder(d = d, h = l);

    // Inner cylinder
    cylinder(d = d - 2 * t, h = l);
}