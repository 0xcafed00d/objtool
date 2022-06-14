package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

func exitOnError(e error, msg string) {
	if e != nil {
		abend(msg + " : " + e.Error())
	}
}

func abend(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(-1)
}

// Config settings from invocation flags
type Config struct {
	help           bool
	outputFilename string
}

var config Config

func init() {

	flag.BoolVar(&config.help, "h", false, "display help")
	flag.StringVar(&config.outputFilename, "o", "", "name of output file. Output written to stdout if omitted")

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "objtool: displays information and modifies a wavefront OBJ 3d model file")
		fmt.Fprintln(os.Stderr, "  Usage: objtool [options] <input file name> ")
		flag.PrintDefaults()
	}
}

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
	t_idn int
	n_idx int
}

type Face struct {
	faceVertex []FaceVertex
}

type ObjFile struct {
	vertices      []Vertex
	vertexNormals []VertexNormal
	texCoords     []TextureCoord
	faces         []Face
	material      string
}

func loadFile(r io.Reader, lineFunc func(string) error) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
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

func main() {
	flag.Parse()

	if len(flag.Args()) == 0 || config.help {
		flag.Usage()
		os.Exit(1)
	}

}
