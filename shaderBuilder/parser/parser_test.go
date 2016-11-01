package parser

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"regexp"
)

func fixNewLines(str string) string {
	re := regexp.MustCompile(`[\n|\r]+`)
	return re.ReplaceAllString(str, "\n")
}

func TestGenerator(t *testing.T) {
	expectedFragFd, err := os.Open("./test/expected.frag")
	assert.NoError(t, err)
	defer expectedFragFd.Close()

	expectedVertFd, err := os.Open("./test/expected.vert")
	assert.NoError(t, err)
	defer expectedVertFd.Close()

	shaderPath := "./test/test.glsl"
	frag := new(bytes.Buffer)
	vert := new(bytes.Buffer)
	err = ParseFile(shaderPath, vert, frag, nil)

	expectedFrag := new(bytes.Buffer)
	expectedFrag.ReadFrom(expectedFragFd)
	assert.Equal(t, fixNewLines(expectedFrag.String()), fixNewLines(frag.String()))

	expectedVert := new(bytes.Buffer)
	expectedVert.ReadFrom(expectedVertFd)
	assert.Equal(t, fixNewLines(expectedVert.String()), fixNewLines(vert.String()))
}
