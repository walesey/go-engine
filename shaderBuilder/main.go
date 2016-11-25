package main

import (
	"bytes"
	"fmt"
	"os"
	"regexp"

	"github.com/walesey/go-engine/shaderBuilder/parser"
)

func main() {
	srcFile := "./index.glsl"
	if len(os.Args) >= 2 {
		srcFile = os.Args[1]
	}

	mode := "vert"
	if len(os.Args) >= 3 {
		mode = os.Args[2]
	}

	out := new(bytes.Buffer)
	switch mode {
	case "vert":
		parser.ParseFile(srcFile, out, nil, nil)
	case "frag":
		parser.ParseFile(srcFile, nil, out, nil)
	case "geo":
		parser.ParseFile(srcFile, nil, nil, out)
	default:
		panic("Invalid shader type: " + mode)
	}
	output := out.String()

	re := regexp.MustCompile("\n[\\s]+\n[\\s]+\n")
	fmt.Println(re.ReplaceAllString(output, "\n\n"))
}
