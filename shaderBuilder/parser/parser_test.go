package parser

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerator(t *testing.T) {
	inFd, err := os.Open("./test/test.glsl")
	assert.NoError(t, err)
	defer inFd.Close()

	expectedFragFd, err := os.Open("./test/expected.frag")
	assert.NoError(t, err)
	defer expectedFragFd.Close()

	expectedVertFd, err := os.Open("./test/expected.vert")
	assert.NoError(t, err)
	defer expectedVertFd.Close()

	frag := new(bytes.Buffer)
	vert := new(bytes.Buffer)
	New(inFd, frag, vert, nil).Parse()

	expectedFrag := new(bytes.Buffer)
	expectedFrag.ReadFrom(expectedFragFd)
	expectedFragStr := strings.Replace(expectedFrag.String(), "\r", "", -1)
	assert.Equal(t, strings.Trim(expectedFragStr, " \n"), strings.Trim(frag.String(), " \n"))

	expectedVert := new(bytes.Buffer)
	expectedVert.ReadFrom(expectedVertFd)
	expectedVertStr := strings.Replace(expectedVert.String(), "\r", "", -1)
	assert.Equal(t, strings.Trim(expectedVertStr, " \n"), strings.Trim(vert.String(), " \n"))
}
