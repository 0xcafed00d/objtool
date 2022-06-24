package main

import (
	"fmt"
	"io"
)

func saveFile(w io.Writer, objFile *ObjFile) error {
	for _, v := range objFile.vertices {
		if v.w == 1.0 {
			fmt.Fprintf(w, "v %f %f %f\n", v.x, v.y, v.z)
		} else {
			fmt.Fprintf(w, "v %f %f %f %f\n", v.x, v.y, v.z, v.w)
		}
	}
	fmt.Fprintln(w, "")

	for _, tc := range objFile.texCoords {
		if tc.w == 0.0 {
			fmt.Fprintf(w, "vt %f %f\n", tc.u, tc.v)
		} else {
			fmt.Fprintf(w, "vt %f %f %f\n", tc.u, tc.v, tc.w)

		}
	}
	fmt.Fprintln(w, "")

	for _, vn := range objFile.vertexNormals {
		fmt.Fprintf(w, "vn %f %f %f\n", vn.x, vn.y, vn.z)
	}
	fmt.Fprintln(w, "")

	for _, f := range objFile.faces {
		fmt.Fprintf(w, "f")
		for _, v := range f.faceVertex {
			fmt.Fprintf(w, " %d", v.v_idx)
			if v.has_t {
				fmt.Fprintf(w, "/%d", v.t_idx)
			}
			if v.has_n {
				fmt.Fprintf(w, "/%d", v.n_idx)
			}
		}
		fmt.Fprintf(w, "\n")
	}

	_, err := fmt.Fprintln(w, "")
	return err
}
