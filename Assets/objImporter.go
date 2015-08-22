package assets

import (
	"strings"
	"strconv"
	"bufio"
	"os"
	"fmt"
	"log"
)

type ObjData struct{
	Name string
	Indicies []uint32
	Vertices []float32
} 

//returns corresponding index (0,1,2...)
func (obj *ObjData) pushVert( x,y,z,nx,ny,nz,u,v float32 ) uint32 {
	obj.Vertices = append(obj.Vertices, x,y,z,nx,ny,nz,u,v )
	return (uint32)((len(obj.Vertices) / 8) - 1)
}

func (obj *ObjData) pushIndex( f1,f2,f3 uint32 ) {
	obj.Indicies = append(obj.Indicies, f1,f2,f3 )
}

func (obj *ObjData) processFaceVertex( token string, vertexList, uvList, normalList []float32 ) uint32 {
	face := strings.Split(token, "/")
	var index int32

	//vertex
	vx := (float32)(0.0)
	vy := (float32)(0.0)
	vz := (float32)(0.0)
	if( face[0] != "" ){
		index = (sti(face[0])-1) * 3
		vx = vertexList[index]
		vy = vertexList[index+1]
		vz = vertexList[index+2]
	}

	//texture
	vtx := (float32)(0.0)
	vty := (float32)(0.0)
	if( face[1] != "" ){
		index = (sti(face[1])-1) * 2
		vtx = uvList[index]
		vty = uvList[index+1]
	}
	//normal
	vnx := (float32)(0.0)
	vny := (float32)(0.0)
	vnz := (float32)(0.0)
	if( face[2] != "" ){
		index = (sti(face[2])-1) * 3
		vnx = normalList[index]
		vny = normalList[index+1]
		vnz = normalList[index+2]
	}

	return obj.pushVert( vx,vy,vz, vnx,vny,vnz, vtx,vty )
}

func (obj *ObjData) processFace( line string, vertexList, uvList, normalList []float32 ){
	tokens := strings.Fields(line)
	if tokens[0] == "f" {
		i1 := obj.processFaceVertex(tokens[1], vertexList, uvList, normalList)
		i2 := obj.processFaceVertex(tokens[2], vertexList, uvList, normalList)
		i3 := obj.processFaceVertex(tokens[3], vertexList, uvList, normalList)
		obj.pushIndex(i1, i2, i3)
	}
}

func ImportObj(filePath string) *ObjData {

	obj := &ObjData{ Indicies: make([]uint32, 0, 0), Vertices: make([]float32, 0, 0) }
	vertexList := make([]float32, 0, 0)
	uvList := make([]float32, 0, 0)
	normalList := make([]float32, 0, 0)

	file, err := os.Open(filePath)
	if err != nil {
	    log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.Fields(line)
		dataType := tokens[0]
		if dataType == "o" { //sub mesh
			obj.Name = tokens[1]
		} else if dataType == "v" { //xyz vertex
			vertexList = append(vertexList, stf(tokens[1]), stf(tokens[2]), stf(tokens[3]) )
		} else if dataType == "vt" { //uv coord
			uvList = append(uvList, stf(tokens[1]), stf(tokens[2]) )
		} else if dataType == "vn" { //xyz vertex normal
			normalList = append(normalList, stf(tokens[1]), stf(tokens[2]), stf(tokens[3]) )
		} else if dataType == "usemtl" { //mtl material
			fmt.Println("NOTE: mtl not yet supported: " + line)
		} else if dataType == "f" { // v/t/n face
			//check triangles
			if len(tokens) != 4 {
				fmt.Println("WARNING: ImportObj only supports triangles: " + line)
			}
			obj.processFace( line, vertexList, uvList, normalList );
		}
	}

	if err := scanner.Err(); err != nil {
	    log.Fatal(err)
	}

	return obj
}

//string to float32
func stf(s string) float32 {
	f, err := strconv.ParseFloat(s, 32)
	if err != nil {
	    log.Fatal(err)
	}
	return (float32)(f)
}

//string to int32
func sti(s string) int32 {
	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
	    log.Fatal(err)
	}
	return (int32)(i)
}