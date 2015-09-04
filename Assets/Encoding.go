package assets

import(
	"fmt"
	"bytes"
	"image"
	"image/png"
	"encoding/json"

	"github.com/Walesey/goEngine/renderer"
)

type MaterialJSON struct {
    LightingMode int32 `json:"lightingMode"`
    Diffuse []byte `json:"diffuse"`
    Normal []byte `json:"normal"`
    Specular []byte `json:"specular"`
    Roughness []byte `json:"roughness"`
}

type GeometryJSON struct {
    Indicies []uint32 `json:"indicies"`
    Verticies []float32 `json:"verticies"`
    CullBackface bool `json:"cullBackface"`
}

//
func EncodeMaterial( material *renderer.Material ) string {
	matJSON := MaterialJSON{ 
		material.LightingMode, 
		EncodeImage( material.Diffuse ), 
		EncodeImage( material.Normal ), 
		EncodeImage( material.Specular ), 
		EncodeImage( material.Roughness ), 
	}
	data, err := json.Marshal(matJSON)
	if err != nil {
		panic(err)
	}
	return string(data)
}

//
func DecodeMaterial( input string ) *renderer.Material{
	var matJSON MaterialJSON
	err := json.Unmarshal([]byte(input), &matJSON)
	if err != nil {
		panic(err)
	}
	mat := renderer.CreateMaterial()
    mat.Diffuse = DecodeImage( matJSON.Diffuse )
    mat.Normal = DecodeImage( matJSON.Normal )
    mat.Specular = DecodeImage( matJSON.Specular )
    mat.Roughness = DecodeImage( matJSON.Roughness )
    mat.LightingMode = matJSON.LightingMode
    return mat
}

//
func EncodeGeometry( geometry *renderer.Geometry ) string {
	geomJSON := GeometryJSON{ geometry.Indicies, geometry.Verticies, geometry.CullBackface }
	data, err := json.Marshal(geomJSON)
	if err != nil {
		panic(err)
	}
	return string(data)
}

//
func DecodeGeometry( input string ) *renderer.Geometry{
	var geomJSON GeometryJSON
	err := json.Unmarshal([]byte(input), &geomJSON)
	if err != nil {
		panic(err)
	}
	geometry := renderer.CreateGeometry( geomJSON.Indicies, geomJSON.Verticies )
	geometry.CullBackface = geomJSON.CullBackface
	return geometry
}

func EncodeImage(img image.Image) []byte {
	if img == nil {
		return []byte{}
	}
	b := make([]byte, 0, 0)
	buf := bytes.NewBuffer(b)
	err := png.Encode(buf, img)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func DecodeImage(pngImg []byte) image.Image {
	b := bytes.NewBuffer(pngImg)
	img,err := png.Decode(b)
	if err != nil {
		fmt.Println(err)
	}
	return img
} 