package main

import (
	"bytes"
	"fmt"
	"image/color"
	"os"
	"runtime"

	"github.com/walesey/go-engine/assets"
	"github.com/walesey/go-engine/controller"
	"github.com/walesey/go-engine/engine"
	"github.com/walesey/go-engine/glfwController"
	"github.com/walesey/go-engine/networking"
	"github.com/walesey/go-engine/opengl"
	"github.com/walesey/go-engine/renderer"
	"github.com/walesey/go-engine/util"

	vmath "github.com/walesey/go-engine/vectormath"
)

type Player struct {
	clientId string
	node     *renderer.Node
	velocity vmath.Vector3
}

// Player is an object that moves around with a certain velocity
func (p *Player) Update(dt float64) {
	p.node.SetTranslation(p.node.Translation.Add(p.velocity.MultiplyScalar(dt)))
}

func init() {
	// Use all cpu cores
	runtime.GOMAXPROCS(runtime.NumCPU())
	//Set default glfw controller
	controller.SetDefaultConstructor(glfwController.NewActionMap)
}

//
func main() {
	server := len(os.Args) > 1 && os.Args[1] == "server"
	var gameEngine engine.Engine
	var glRenderer *opengl.OpenglRenderer
	network := networking.NewNetwork()

	// Start server or connect to server
	if server {
		gameEngine = engine.NewHeadlessEngine()
		network.StartServer(1234)
	} else {
		glRenderer = &opengl.OpenglRenderer{
			WindowTitle:  "Networking example",
			WindowWidth:  800,
			WindowHeight: 800,
		}
		gameEngine = engine.NewEngine(glRenderer)
		network.ConnectClient("127.0.0.1:1234")
	}
	gameEngine.AddUpdatable(network)

	// map containing each player's entity
	players := make(map[string]*Player)

	//Networked Game events
	network.ClientJoinedEvent(func(clientId string) {
		fmt.Println("client joined, clientId: ", clientId)
		network.TriggerOnServerAndClients("spawn", []byte(clientId))
		for _, player := range players {
			network.TriggerEvent("spawn", clientId, []byte(player.clientId))
		}
	})

	network.RegisterEvent("spawn", func(clientId string, data []byte) {
		player := &Player{clientId: string(data)}
		players[clientId] = player
		if network.IsClient() {
			gameEngine.AddUpdatable(player)
			player.node = renderer.CreateNode()
			boxGeometry := renderer.CreateBox(1, 1)
			boxGeometry.Material = renderer.CreateMaterial()
			boxGeometry.SetColor(color.NRGBA{254, 0, 0, 254})
			boxGeometry.CullBackface = false
			player.node.SetTranslation(vmath.Vector3{X: -20})
			player.node.SetOrientation(vmath.AngleAxis(1.57, vmath.Vector3{Y: 1}))
			player.node.Add(boxGeometry)
			gameEngine.AddSpatial(player.node)
		}
	})

	network.RegisterEvent("move", func(clientId string, data []byte) {
		buf := bytes.NewBuffer(data)
		position := util.Vector3frombytes(buf)
		velocity := util.Vector3frombytes(buf)
		if player, ok := players[clientId]; ok {
			player.node.SetTranslation(position)
			player.velocity = velocity
		}
	})

	// client setup
	gameEngine.Start(func() {
		if network.IsClient() {

			//lighting
			glRenderer.CreateLight(
				0.1, 0.1, 0.1, //ambient
				0.5, 0.5, 0.5, //diffuse
				0.7, 0.7, 0.7, //specular
				true, vmath.Vector3{0.7, 0.7, 0.7}, //direction
				0, //index
			)

			// Sky cubemap
			skyImg, err := assets.ImportImage("resources/cubemap.png")
			if err == nil {
				gameEngine.Sky(assets.CreateMaterial(skyImg, nil, nil, nil), 999999)
			}

			// input/controller manager
			controllerManager := glfwController.NewControllerManager(glRenderer.Window)

			// networked movement controls
			move := func(velocity vmath.Vector3) {
				network.TriggerOnServerAndClients("move", util.SerializeArgs(network.ClientToken(), velocity))
			}

			customController := controller.CreateController()
			controllerManager.AddController(customController.(glfwController.Controller))
			customController.BindAction(func() { move(vmath.Vector3{Y: 1}) }, controller.KeyW, controller.Press)
			customController.BindAction(func() { move(vmath.Vector3{Z: 1}) }, controller.KeyA, controller.Press)
			customController.BindAction(func() { move(vmath.Vector3{Y: -1}) }, controller.KeyS, controller.Press)
			customController.BindAction(func() { move(vmath.Vector3{Z: -1}) }, controller.KeyD, controller.Press)
		}
	})
}
