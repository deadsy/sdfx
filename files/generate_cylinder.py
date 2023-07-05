import math


# Cylinder is shifted by {+radius} along the Z axis to be above the XY plane.
def generate_cylinder_stl(radius, length, num_triangles, output_file):
    with open(output_file, "w") as stl_file:
        stl_file.write("solid cylinder\n")

        angle_step = 2 * math.pi / num_triangles

        for i in range(num_triangles):
            angle1 = i * angle_step
            angle2 = (i + 1) * angle_step

            y1, z1 = radius * math.cos(angle1), radius * math.sin(angle1)
            y2, z2 = radius * math.cos(angle2), radius * math.sin(angle2)

            # Base triangle
            stl_file.write("facet normal -1.0 0.0 0.0\n")
            stl_file.write("outer loop\n")
            stl_file.write(f"vertex 0.0 0.0 {0.0+radius}\n")
            stl_file.write(f"vertex 0.0 {y2} {z2+radius}\n")
            stl_file.write(f"vertex 0.0 {y1} {z1+radius}\n")
            stl_file.write("endloop\n")
            stl_file.write("endfacet\n")

            # Top triangle
            stl_file.write("facet normal 1.0 0.0 0.0\n")
            stl_file.write("outer loop\n")
            stl_file.write(f"vertex {length} 0.0 {0.0+radius}\n")
            stl_file.write(f"vertex {length} {y1} {z1+radius}\n")
            stl_file.write(f"vertex {length} {y2} {z2+radius}\n")
            stl_file.write("endloop\n")
            stl_file.write("endfacet\n")

            # Side triangles
            normal = f"0.0 {math.cos(angle1)} {math.sin(angle1)}"
            stl_file.write(f"facet normal {normal}\n")
            stl_file.write("outer loop\n")
            stl_file.write(f"vertex 0.0 {y1} {z1+radius}\n")
            stl_file.write(f"vertex 0.0 {y2} {z2+radius}\n")
            stl_file.write(f"vertex {length} {y1} {z1+radius}\n")
            stl_file.write("endloop\n")
            stl_file.write("endfacet\n")

            normal = f"0.0 {math.cos(angle2)} {math.sin(angle2)}"
            stl_file.write(f"facet normal {normal}\n")
            stl_file.write("outer loop\n")
            stl_file.write(f"vertex 0.0 {y2} {z2+radius}\n")
            stl_file.write(f"vertex {length} {y2} {z2+radius}\n")
            stl_file.write(f"vertex {length} {y1} {z1+radius}\n")
            stl_file.write("endloop\n")
            stl_file.write("endfacet\n")

        stl_file.write("endsolid cylinder\n")


output_filename = "cylinder.stl"
radius = 12
length = 100
num_triangles = 36
generate_cylinder_stl(radius, length, num_triangles, output_filename)
