package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/walesey/go-engine/assets"
	"github.com/walesey/go-engine/examples"

	"github.com/codegangsta/cli"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
	// Use all cpu cores
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	app := cli.NewApp()
	app.Name = "goEngine"
	app.Usage = ""
	app.EnableBashCompletion = true
	app.Commands = []cli.Command{
		{
			Name:   "material",
			Usage:  "Import material from file and save to a .asset file",
			Action: materialImport,
		},
		{
			Name:   "image",
			Usage:  "Import texture from file and save to a .asset file",
			Action: imageImport,
		},
		{
			Name:   "geometry",
			Usage:  "Import obj geometry file and save to a .asset file",
			Action: geometryImport,
		},
		{
			Name:   "list",
			Usage:  "List all assets in a .asset file",
			Action: list,
		},
		{
			Name:   "remove",
			Usage:  "Remove an asset from a .asset file",
			Action: remove,
		},

		//DEMOS
		{
			Name:   "demo",
			Usage:  "very basic demo",
			Action: examples.Demo,
		},
		{
			Name:   "particles",
			Usage:  "run a particle effect example",
			Action: examples.Particles,
		},
		{
			Name:   "gun",
			Usage:  "run a demo of a gun model",
			Action: examples.GunDemo,
		},
		{
			Name:   "physics",
			Usage:  "run the physics demo",
			Action: examples.PhysicsDemo,
		},
	}

	app.Run(os.Args)
}

//CLI remove asset from file
func remove(c *cli.Context) {
	if len(c.Args()) != 2 {
		fmt.Println("Usage: goEngine material <assetFile> <name> ")
		return
	}
	assetLib, _ := assets.LoadAssetLibrary(c.Args()[0])
	delete(assetLib.Assets, c.Args()[1])
	assetLib.SaveToFile(c.Args()[0])
}

//CLI list assets from file
func list(c *cli.Context) {
	if len(c.Args()) != 1 {
		fmt.Println("Usage: goEngine material <assetFile> ")
		return
	}
	assetLib, _ := assets.LoadAssetLibrary(c.Args()[0])
	for name := range assetLib.Assets {
		fmt.Println(name, ": ", assetLib.Assets[name].Type)
	}
}

//CLI material creator
func materialImport(c *cli.Context) {
	if len(c.Args()) != 6 {
		fmt.Println("Usage: goEngine material <assetFile> <name> <albedoFile> <normalFile> <specFile> <roughnessFile> ")
		return
	}
	diffuseMap := assets.ImportImage(c.Args()[2])
	normalMap := assets.ImportImage(c.Args()[3])
	specMap := assets.ImportImage(c.Args()[4])
	roughnessMap := assets.ImportImage(c.Args()[5])
	mat := assets.CreateMaterial(diffuseMap, normalMap, specMap, roughnessMap)
	assetLib, _ := assets.LoadAssetLibrary(c.Args()[0])
	assetLib.AddMaterial(c.Args()[1], mat)
	assetLib.SaveToFile(c.Args()[0])
}

//CLI geometry creator
func geometryImport(c *cli.Context) {
	if len(c.Args()) != 3 {
		fmt.Println("Usage: goEngine geometry <assetFile> <name> <objFile> ")
		return
	}
	geometry := assets.ImportObj(c.Args()[2])
	assetLib, _ := assets.LoadAssetLibrary(c.Args()[0])
	assetLib.AddGeometry(c.Args()[1], geometry)
	assetLib.AddMaterial(fmt.Sprint(c.Args()[1], "Mat"), *geometry.Material)
	assetLib.SaveToFile(c.Args()[0])
}

//CLI image creator
func imageImport(c *cli.Context) {
	if len(c.Args()) != 3 {
		fmt.Println("Usage: goEngine image <assetFile> <name> <imageFile>")
		return
	}
	imageAsset := assets.ImportImage(c.Args()[2])
	assetLib, _ := assets.LoadAssetLibrary(c.Args()[0])
	assetLib.AddImage(c.Args()[1], imageAsset)
	assetLib.SaveToFile(c.Args()[0])
}
