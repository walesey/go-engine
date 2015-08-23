package assets

import (
	"strings"
	"strconv"
	"bufio"
	"os"
	"log"
)

//Vertices format : x,y,z,   nx,ny,nz,   u,v 
//Indicies format : f1,f2,f3 (triangles)
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

func (obj *ObjData) pushIndex( indicies ...uint32 ) {
	obj.Indicies = append(obj.Indicies, indicies... )
}

//parses a single triangle vertex, returning the newly generated index
func (obj *ObjData) processFaceVertex( token string, vertexList, uvList, normalList []float32 ) uint32 {
	face := strings.Split(token, "/")
	var index int32
	
	//vertex
	vx := (float32)(0.0)
	vy := (float32)(0.0)
	vz := (float32)(0.0)
	if len(face) > 0 && face[0] != "" {
		index = (sti(face[0])-1) * 3
		vx = vertexList[index]
		vy = vertexList[index+1]
		vz = vertexList[index+2]
	}
	
	//texture
	vtx := (float32)(0.0)
	vty := (float32)(0.0)
	if len(face) > 1 && face[1] != "" {
		index = (sti(face[1])-1) * 2
		vtx = uvList[index]
		vty = uvList[index+1]
	}
	
	//normal
	vnx := (float32)(0.0)
	vny := (float32)(0.0)
	vnz := (float32)(0.0)
	if len(face) > 2 && face[2] != "" {
		index = (sti(face[2])-1) * 3
		vnx = normalList[index]
		vny = normalList[index+1]
		vnz = normalList[index+2]
	}
	
	return obj.pushVert( vx,vy,vz, vnx,vny,vnz, vtx,vty )
}

//Processes a polygonal face by splitting it into triangles 
func (obj *ObjData) processFace( line string, vertexList, uvList, normalList []float32 ){
	tokens := strings.Fields(line)
	if tokens[0] == "f" {
		tokens = append(tokens[:0], tokens[1:]...)
		for len(tokens) > 0 {
			tempTokens := make([]string, 0, 0)
			for i:=0 ; i < (len(tokens)-1); i+=2 {
				obj.pushIndex( obj.processFaceVertex(tokens[i], vertexList, uvList, normalList) )
				obj.pushIndex( obj.processFaceVertex(tokens[i+1], vertexList, uvList, normalList) )
				if len(tokens) > (i+2){
					obj.pushIndex( obj.processFaceVertex(tokens[i+2], vertexList, uvList, normalList) )
				} else {
					obj.pushIndex( obj.processFaceVertex(tokens[0], vertexList, uvList, normalList) )
				}
				if len(tokens) > 4 {
					tempTokens = append(tempTokens, tokens[i])
				}
			}
			if len(tokens) > 4 && len(tokens)%2 == 1 {
				tempTokens = append(tempTokens, tokens[len(tokens)-1])
			}
			tokens = tempTokens
		}
	}
}

//imports an obj filePath into an ObjData reference containing index and vertex buffers
func ImportObj(filePath string) (*ObjData, error) {

	obj := &ObjData{ Indicies: make([]uint32, 0, 0), Vertices: make([]float32, 0, 0) }
	vertexList := make([]float32, 0, 0)
	uvList := make([]float32, 0, 0)
	normalList := make([]float32, 0, 0)

	file, err := os.Open(filePath)
	if err != nil {
		return obj, err
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
			//TODO
		} else if dataType == "f" { // v/t/n face
			obj.processFace( line, vertexList, uvList, normalList );
		}
	}

	if err := scanner.Err(); err != nil {
		return obj, err
	}

	return obj, nil
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