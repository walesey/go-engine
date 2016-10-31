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

	src, err := os.Open(srcFile)
	if err != nil {
		panic(err)
	}

	frag := new(bytes.Buffer)
	vert := new(bytes.Buffer)
	parser.New(src, frag, vert, nil).Parse()
	fmt.Print("#frag")
	fmt.Print(frag.String())
	fmt.Print("#vert")
	fmt.Print(vert.String())
}
