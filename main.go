package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
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
	materialName   string
	centerOrigin   bool
	centerOriginX  bool
	centerOriginY  bool
	centerOriginZ  bool
	y0align        bool
	triangle       bool
	scale          float64
	resizeX        float64
	resizeY        float64
	resizeZ        float64
}

var config Config

// a more condensed printout of the program options than the default go lib version
func PrintOptions() {
	flag.VisitAll(func(f *flag.Flag) {
		var b strings.Builder
		fmt.Fprintf(&b, "  -") // Two spaces before -; see next two comments.

		if len(f.Name) > 0 {
			b.WriteString(f.Name)
		}

		b.WriteString("\t")
		b.WriteString(strings.ReplaceAll(f.Usage, "\n", "\n    \t"))

		fmt.Fprint(os.Stderr, b.String(), "\n")
	})
}

func init() {

	flag.BoolVar(&config.help, "h", false, "display help")
	flag.StringVar(&config.outputFilename, "o", "", "name of output file. Output written to stdout if omitted")
	flag.StringVar(&config.materialName, "m", "", "sets the model to use the specified material")
	flag.BoolVar(&config.centerOrigin, "c", false, "Move the object so its center is (0,0,0)")
	flag.BoolVar(&config.centerOriginX, "cx", false, "Move the object along the X axis so its center is X=0")
	flag.BoolVar(&config.centerOriginY, "cy", false, "Move the object along the Y axis so its center is Y=0")
	flag.BoolVar(&config.centerOriginZ, "cz", false, "Move the object along the Z axis so its center is Z=0")
	flag.BoolVar(&config.y0align, "y0", false, "align the base of object to Y=0 (the ground)")
	flag.BoolVar(&config.triangle, "t", false, "convert all faces to triangles")
	flag.Float64Var(&config.scale, "s", 100.0, "scale the object to s%")
	flag.Float64Var(&config.resizeX, "rx", -1.0, "resize object along x to specified size, all other axis are scale in proportion")
	flag.Float64Var(&config.resizeY, "ry", -1.0, "resize object along y to specified size, all other axis are scale in proportion")
	flag.Float64Var(&config.resizeZ, "rz", -1.0, "resize object along z to specified size, all other axis are scale in proportion")

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "objtool: displays information and modifies a wavefront OBJ 3d model file")
		fmt.Fprintln(os.Stderr, "  Usage: objtool [options] <input file name> ")
		PrintOptions()
	}
}

func main() {
	flag.Parse()

	if len(flag.Args()) == 0 || config.help {
		flag.Usage()
		os.Exit(1)
	}

	input := flag.Args()[0]

	infile, err := os.Open(input)
	exitOnError(err, "Cant Open Inputfile")
	defer infile.Close()

	objFile := ObjFile{}

	err = loadFile(infile, func(line string) error {
		return processLine(line, &objFile)
	})

	info := getInfo(&objFile)
	displayInfo(&info)

	processOptions(&config, &objFile)

	postinfo := getInfo(&objFile)
	displayInfo(&postinfo)

	exitOnError(err, "Error Reading Inputfile")
}
