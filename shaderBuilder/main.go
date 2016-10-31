package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/walesey/go-engine/shaderBuilder/parser"
)

func main() {
	srcFile := "./index.glsl"
	if len(os.Args) >= 2 {
		srcFile = os.Args[1]
	}

	vert, frag := new(bytes.Buffer),new(bytes.Buffer)
	parser.ParseFile(srcFile, frag, vert, nil)
	fmt.Println("#vert")
	fmt.Println(vert.String())
	fmt.Println("#frag")
	fmt.Println(frag.String())
}
