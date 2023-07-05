def write_stl(filename):
    with open(filename, "w") as f:
        f.write("solid cube\n")
        for face in faces:
            f.write("  facet normal 0 0 0\n")
            f.write("    outer loop\n")
            for vertex in face:
                f.write("      vertex {} {} {}\n".format(*vertex))
            f.write("    endloop\n")
            f.write("  endfacet\n")
        f.write("endsolid cube")


width = 12
height = 22
length = 100

# Define the 8 vertices of the cube
vertices = [
    (0, 0, 0),
    (length, 0, 0),
    (length, width, 0),
    (0, width, 0),
    (0, 0, height),
    (length, 0, height),
    (length, width, height),
    (0, width, height),
]

# Define the 12 triangles composing the cube
faces = [
    (vertices[0], vertices[1], vertices[2]),
    (vertices[0], vertices[2], vertices[3]),
    (vertices[0], vertices[1], vertices[5]),
    (vertices[0], vertices[5], vertices[4]),
    (vertices[1], vertices[2], vertices[6]),
    (vertices[1], vertices[6], vertices[5]),
    (vertices[2], vertices[3], vertices[7]),
    (vertices[2], vertices[7], vertices[6]),
    (vertices[3], vertices[0], vertices[4]),
    (vertices[3], vertices[4], vertices[7]),
    (vertices[4], vertices[5], vertices[6]),
    (vertices[4], vertices[6], vertices[7]),
]

write_stl("cube.stl")
