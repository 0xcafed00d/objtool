package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
)

type Vertex struct {
	x float64
	y float64
	z float64
	w float64
}

type VertexNormal struct {
	x float64
	y float64
	z float64
}

type TextureCoord struct {
	u float64
	v float64
	w float64
}

type FaceVertex struct {
	v_idx int
	t_idx int
	has_t bool
	n_idx int
	has_n bool
}

type Face struct {
	faceVertex []FaceVertex
}

type ObjFile struct {
	vertices      []Vertex
	vertexNormals []VertexNormal
	texCoords     []TextureCoord
	faces         []Face
	materialName  string
}

func loadFile(r io.Reader, lineFunc func(string) error) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		//fmt.Println(line)
		err := lineFunc(line)
		if err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

// turn a space separated list of numbers to a slice of doubles
func string2numbers(numberList string) ([]float64, error) {
	numberList = strings.Trim(numberList, " \t")
	numberStrs := strings.Split(numberList, " ")

	numbers := []float64{}
	for _, n := range numberStrs {
		num, err := strconv.ParseFloat(n, 64)
		if err != nil {
			return nil, err
		}
		numbers = append(numbers, num)
	}
	return numbers, nil
}

func processLineVN(line string, objFile *ObjFile) error {
	n, err := string2numbers(line[2:])
	if err != nil {
		return err
	}
	objFile.vertexNormals = append(objFile.vertexNormals, VertexNormal{n[0], n[1], n[2]})
	return nil
}

func processLineVT(line string, objFile *ObjFile) error {
	n, err := string2numbers(line[2:])
	if err != nil {
		return err
	}
	objFile.texCoords = append(objFile.texCoords, TextureCoord{n[0], n[1], 0})
	return nil
}

func processLineV(line string, objFile *ObjFile) error {
	n, err := string2numbers(line[1:])
	if err != nil {
		return err
	}
	objFile.vertices = append(objFile.vertices, Vertex{n[0], n[1], n[2], 1.0})
	return nil
}

func processLineF(line string, objFile *ObjFile) error {
	line = strings.Trim(line[1:], " \t")
	indexStrs := strings.Split(line, " ")

	indices := []FaceVertex{}
	for _, i := range indexStrs {
		if strings.Contains(i, "/") {
			var err error

			idxs := strings.Split(i, "/")
			fv := FaceVertex{}

			fv.v_idx, err = strconv.Atoi(idxs[0])
			if err != nil {
				return err
			}

			if idxs[1] != "" {
				fv.t_idx, err = strconv.Atoi(idxs[1])
				fv.has_t = true
				if err != nil {
					return err
				}
			}

			if len(idxs) > 2 {
				fv.n_idx, err = strconv.Atoi(idxs[2])
				fv.has_n = true
				if err != nil {
					return err
				}
			}

			indices = append(indices, fv)
		} else {
			idx, err := strconv.Atoi(i)
			if err != nil {
				return err
			}
			indices = append(indices, FaceVertex{v_idx: idx})
		}
	}
	objFile.faces = append(objFile.faces, Face{indices})

	return nil
}

func processLine(line string, objFile *ObjFile) error {
	if len(line) > 0 {
		if strings.HasPrefix(line, "vn") {
			return processLineVN(line, objFile)
		} else if strings.HasPrefix(line, "vt") {
			return processLineVT(line, objFile)
		} else if strings.HasPrefix(line, "v") {
			return processLineV(line, objFile)
		} else if strings.HasPrefix(line, "f") {
			return processLineF(line, objFile)
		}
	}
	return nil
}

type ObjFileInfo struct {
	BBoxMin     Vertex
	BBoxMax     Vertex
	BBoxSize    Vertex
	Origin      Vertex
	FaceCount   int
	FaceTypes   map[int]int
	VertexCount int
}

func getInfo(objFile *ObjFile) ObjFileInfo {

	BBmin := Vertex{math.Inf(1), math.Inf(1), math.Inf(1), 0}
	BBmax := Vertex{math.Inf(-1), math.Inf(-1), math.Inf(-1), 0}

	for _, v := range objFile.vertices {
		BBmin.x = math.Min(v.x, BBmin.x)
		BBmin.y = math.Min(v.y, BBmin.y)
		BBmin.z = math.Min(v.z, BBmin.z)
		BBmax.x = math.Max(v.x, BBmax.x)
		BBmax.y = math.Max(v.y, BBmax.y)
		BBmax.z = math.Max(v.z, BBmax.z)
	}

	info := ObjFileInfo{FaceTypes: map[int]int{}}
	info.VertexCount = len(objFile.vertices)
	info.FaceCount = len(objFile.faces)
	for _, f := range objFile.faces {
		info.FaceTypes[len(f.faceVertex)] += 1
	}

	info.BBoxMax = BBmax
	info.BBoxMin = BBmin
	info.BBoxSize = Vertex{BBmax.x - BBmin.x, BBmax.y - BBmin.y, BBmax.z - BBmin.z, 0.0}
	info.Origin = Vertex{(BBmax.x + BBmin.x) / 2.0, (BBmax.y + BBmin.y) / 2.0, (BBmax.z + BBmin.z) / 2.0, 0.0}

	return info
}

func displayInfo(nfo *ObjFileInfo) {

	fmt.Printf("#Vertex Count: %d \n", nfo.VertexCount)
	fmt.Printf("#  Face Count: %d \n", nfo.FaceCount)
	fmt.Printf("#        Size: {x:%0.6f, y:%0.6f, z:%0.6f} \n", nfo.BBoxSize.x, nfo.BBoxSize.y, nfo.BBoxSize.z)
	fmt.Printf("#      Origin: {x:%0.6f, y:%0.6f, z:%0.6f} \n", nfo.Origin.x, nfo.Origin.y, nfo.Origin.z)
	fmt.Printf("#      Extent: x: %0.6f -> %0.6f \n", nfo.BBoxMin.x, nfo.BBoxMax.x)
	fmt.Printf("#           y: %0.6f -> %0.6f \n", nfo.BBoxMin.y, nfo.BBoxMax.y)
	fmt.Printf("#           z: %0.6f -> %0.6f \n", nfo.BBoxMin.z, nfo.BBoxMax.z)

	fmt.Printf("#   Face Info--->\n")
	for vertices, count := range nfo.FaceTypes {
		fmt.Printf("# Vertex Count: %d Count: %d\n", vertices, count)
	}
	fmt.Printf("\n")

}
