package main

import (
	"fmt"
	"io"
)

func saveFile(w io.Writer, objFile *ObjFile) error {
	for _, v := range objFile.vertices {
		fmt.Fprintf(w, "v %f %f %f\n", v.x, v.y, v.z)
	}
	for _, vn := range objFile.vertexNormals {
		fmt.Fprintf(w, "vn %f %f %f\n", vn.x, vn.y, vn.z)
	}
	for _, tc := range objFile.texCoords {
		fmt.Fprintf(w, "vt %f %f %f\n", tc.u, tc.v, tc.w)
	}
	for _, f := range objFile.faces {
		fmt.Fprintf(w, "f")
		for _, v := range f.faceVertex {
			fmt.Fprintf(w, " %d", v.v_idx+1)
			if v.has_t {
				fmt.Fprintf(w, "/%d", v.t_idx+1)
			}
			if v.has_n {
				fmt.Fprintf(w, "/%d", v.n_idx+1)
			}
		}
		fmt.Fprintf(w, "\n")
	}

	return nil
}
