package main

func convertToTriangles(faces []Face) []Face {
	triangles := []Face{}
	for _, f := range faces {
		if len(f.faceVertex) == 3 {
			triangles = append(triangles, f)
		} else if len(f.faceVertex) > 3 {
			for i := 0; i < len(f.faceVertex)-2; i++ {
				triangles = append(triangles, Face{f.faceVertex[i : i+3]})
			}
		}
	}
	return triangles
}

func translateVertices(objFile *ObjFile, x, y, z float64) {
	for i := range objFile.vertices {
		objFile.vertices[i].x += x
		objFile.vertices[i].y += y
		objFile.vertices[i].z += z
	}
}

func scaleVertices(objFile *ObjFile, s float64) {
	for i := range objFile.vertices {
		objFile.vertices[i].x *= s
		objFile.vertices[i].y *= s
		objFile.vertices[i].z *= s
	}
}

func processOptions(conf *Config, objFile *ObjFile) {
	info := getInfo(objFile)

	if config.materialName != "" {
		objFile.materialName = config.materialName
	}
	if config.triangle {
		objFile.faces = convertToTriangles(objFile.faces)
	}
	if config.y0align {
		translateVertices(objFile, 0, -info.BBoxMin.y, 0)
	}
}
