package assets

import(
	"io/ioutil"
	"encoding/json"
    "compress/gzip"
    "bytes"
	"image"

	"github.com/Walesey/goEngine/renderer"
)

type Asset struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Data string `json:"data"`
}

type AssetLibrary struct {
	Assets map[string]Asset `json:"assets"`
}

func CreateAssetLibrary() *AssetLibrary{
	return &AssetLibrary{ make(map[string]Asset) }
}

func LoadAssetLibrary(fileName string) (*AssetLibrary, error){
    compressedData, err := ioutil.ReadFile(fileName)
    if err != nil {
    	return CreateAssetLibrary(), err
    }
    buf := bytes.NewBuffer(compressedData)
    r, err := gzip.NewReader(buf)
    if err != nil {
		return CreateAssetLibrary(), err
	}
    data, err := ioutil.ReadAll(r)
    if err != nil {
		return CreateAssetLibrary(), err
	}
    defer r.Close()
    var aLib AssetLibrary
    err = json.Unmarshal(data, &aLib)
	if err != nil {
		return CreateAssetLibrary(), err
	}
	return &aLib, nil
}

//
func (al *AssetLibrary) SaveToFile(fileName string){
	aLib := (*al)
	data, err := json.Marshal(aLib)
	if err != nil {
		panic(err)
	}
    var compressedData bytes.Buffer
	w := gzip.NewWriter(&compressedData)
	w.Write([]byte(data))
	w.Close()
    err = ioutil.WriteFile(fileName, compressedData.Bytes(), 0777)
    if err != nil {
        panic(err)
    }
}

//
func (al *AssetLibrary) AddEncodedAsset(name, assetType, data string){
	al.Assets[name] = Asset{Name:name, Type:assetType, Data:data}
}

//
func (al *AssetLibrary) AddGeometry(name string, geometry renderer.Geometry){
	data := EncodeGeometry(geometry)
	al.AddEncodedAsset( name, "geometry", data )
}

//
func (al *AssetLibrary) AddMaterial(name string, geometry renderer.Material){
	data := EncodeMaterial(geometry)
	al.AddEncodedAsset( name, "material", data )
}

//
func (al *AssetLibrary) AddImage( name string, img image.Image ) {
	data := EncodeImage(img)
	al.AddEncodedAsset( name, "image", data )
}

//
func (al *AssetLibrary) GetAssetType(name string) string{
	return al.Assets[name].Type
}

//
func (al *AssetLibrary) GetGeometry(name string) renderer.Geometry{
	data := al.Assets[name].Data
	return DecodeGeometry(data)
}

//
func (al *AssetLibrary) GetMaterial(name string) renderer.Material{
	data := al.Assets[name].Data
	return DecodeMaterial(data)
}

//
func (al *AssetLibrary) GetImage(name string) image.Image{
	data := al.Assets[name].Data
	return DecodeImage(data)
}
