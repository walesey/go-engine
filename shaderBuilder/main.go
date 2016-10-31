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

	out := new(bytes.Buffer)
	parser.New(src, out).Parse()
	fmt.Print(out.String())
}
