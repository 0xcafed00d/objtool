package main

func convertToTriangles(faces []Face) []Face {
	triangles := []Face{}
	for _, f := range faces {
		for i := 1; i < len(f.faceVertex)-1; i++ {
			f := []FaceVertex{f.faceVertex[0], f.faceVertex[i], f.faceVertex[i+1]}
			triangles = append(triangles, Face{f})
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

	if config.materialName != "" {
		objFile.materialName = config.materialName
	}
	if config.triangle {
		objFile.faces = convertToTriangles(objFile.faces)
	}
	if config.scale != 1.0 {
		info := getInfo(objFile)
		translateVertices(objFile, -info.Origin.x, -info.Origin.y, -info.Origin.z)
		scaleVertices(objFile, config.scale)
		translateVertices(objFile, info.Origin.x, info.Origin.y, info.Origin.z)
	}
	if config.centerOrigin {
		info := getInfo(objFile)
		translateVertices(objFile, -info.Origin.x, -info.Origin.y, -info.Origin.z)
	}
	if config.centerOriginX {
		info := getInfo(objFile)
		translateVertices(objFile, -info.Origin.x, 0, 0)
	}
	if config.centerOriginY {
		info := getInfo(objFile)
		translateVertices(objFile, 0, -info.Origin.y, 0)
	}
	if config.centerOriginZ {
		info := getInfo(objFile)
		translateVertices(objFile, 0, 0, -info.Origin.z)
	}
	if config.resizeX != 0 {
		info := getInfo(objFile)
		translateVertices(objFile, -info.Origin.x, -info.Origin.y, -info.Origin.z)
		scaleVertices(objFile, config.resizeX/info.BBoxSize.x)
		translateVertices(objFile, info.Origin.x, info.Origin.y, info.Origin.z)
	}
	if config.resizeY != 0 {
		info := getInfo(objFile)
		translateVertices(objFile, -info.Origin.x, -info.Origin.y, -info.Origin.z)
		scaleVertices(objFile, config.resizeY/info.BBoxSize.y)
		translateVertices(objFile, info.Origin.x, info.Origin.y, info.Origin.z)
	}
	if config.resizeZ != 0 {
		info := getInfo(objFile)
		translateVertices(objFile, -info.Origin.x, -info.Origin.y, -info.Origin.z)
		scaleVertices(objFile, config.resizeZ/info.BBoxSize.z)
		translateVertices(objFile, info.Origin.x, info.Origin.y, info.Origin.z)
	}
	if config.y0align {
		info := getInfo(objFile)
		translateVertices(objFile, 0, -info.BBoxMin.y, 0)
	}
}
