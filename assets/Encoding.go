package assets

import(
	"fmt"
	"bytes"
	"image"
	"image/png"
	"encoding/json"

	"github.com/walesey/go-engine/renderer"
)

type MaterialJSON struct {
    LightingMode int32 `json:"lightingMode"`
    Diffuse string `json:"diffuse"`
    Normal string `json:"normal"`
    Specular string `json:"specular"`
    Roughness string `json:"roughness"`
}

type ImageJson struct {
    ImageData []byte `json:"imageData"`
}

type GeometryJSON struct {
    Indicies []uint32 `json:"indicies"`
    Verticies []float32 `json:"verticies"`
    CullBackface bool `json:"cullBackface"`
}

//
func EncodeMaterial( material renderer.Material ) string {
	matJSON := MaterialJSON{ LightingMode: material.LightingMode }
	if material.Diffuse != nil {
		matJSON.Diffuse = EncodeImage( material.Diffuse )
	}
	if material.Normal != nil {
		matJSON.Normal = EncodeImage( material.Normal )
	}
	if material.Specular != nil {
		matJSON.Specular = EncodeImage( material.Specular )
	}
	if material.Roughness != nil {
		matJSON.Roughness = EncodeImage( material.Roughness )
	}
	data, err := json.Marshal(matJSON)
	if err != nil {
		panic(err)
	}
	return string(data)
}

//
func DecodeMaterial( input string ) renderer.Material{
	var matJSON MaterialJSON
	err := json.Unmarshal([]byte(input), &matJSON)
	if err != nil {
		panic(err)
	}
	mat := renderer.CreateMaterial()
	if matJSON.Diffuse != "" {
 	   mat.Diffuse = DecodeImage( matJSON.Diffuse )
	}
	if matJSON.Normal != "" {
    	mat.Normal = DecodeImage( matJSON.Normal )
	}
	if matJSON.Specular != "" {
    	mat.Specular = DecodeImage( matJSON.Specular )
	}
	if matJSON.Roughness != "" {
    	mat.Roughness = DecodeImage( matJSON.Roughness )
	}
    mat.LightingMode = matJSON.LightingMode
    return mat
}

//
func EncodeGeometry( geometry renderer.Geometry ) string {
	geomJSON := GeometryJSON{ geometry.Indicies, geometry.Verticies, geometry.CullBackface }
	data, err := json.Marshal(geomJSON)
	if err != nil {
		panic(err)
	}
	return string(data)
}

//
func DecodeGeometry( input string ) renderer.Geometry{
	var geomJSON GeometryJSON
	err := json.Unmarshal([]byte(input), &geomJSON)
	if err != nil {
		panic(err)
	}
	geometry := renderer.CreateGeometry( geomJSON.Indicies, geomJSON.Verticies )
	geometry.CullBackface = geomJSON.CullBackface
	return geometry
}

//
func EncodeImage(img image.Image) string {
	b := make([]byte, 0, 0)
	buf := bytes.NewBuffer(b)
	err := png.Encode(buf, img)
	if err != nil {
		panic(err)
	}
	data, err := json.Marshal(ImageJson{ buf.Bytes() })
	if err != nil {
		panic(err)
	}
	return string(data)
}

//
func DecodeImage(pngImg string) image.Image {
	var imageJSON ImageJson
	err := json.Unmarshal([]byte(pngImg), &imageJSON)
	if err != nil {
		panic(err)
	}
	b := bytes.NewBuffer(imageJSON.ImageData)
	img,err := png.Decode(b)
	if err != nil {
		fmt.Println(err)
	}
	return img
} 