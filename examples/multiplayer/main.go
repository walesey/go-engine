package main

import (
	"bytes"
	"fmt"
	"image/color"
	"os"
	"runtime"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/walesey/go-engine/assets"
	"github.com/walesey/go-engine/controller"
	"github.com/walesey/go-engine/engine"
	"github.com/walesey/go-engine/glfwController"
	"github.com/walesey/go-engine/networking"
	"github.com/walesey/go-engine/opengl"
	"github.com/walesey/go-engine/renderer"
	"github.com/walesey/go-engine/util"
)

/*
	This is an example of a client/server multiplayer game
	To run the server use: go run main.go server
	To run a client use: go run main.go
*/

var serverPort = 1234
var serverAddr = "127.0.0.1"
var startingPosition = mgl32.Vec2{20, 20}

type Player struct {
	clientId string // id of the client associated with this player
	node     *renderer.Node
	position mgl32.Vec2
	velocity mgl32.Vec2
}

// Update the players position based on it's velocity
func (p *Player) Update(dt float64) {
	p.position = p.position.Add(p.velocity.Mul(float32(dt)))
	p.node.SetTranslation(p.position.Vec3(0))
}

func init() {
	// Use all cpu cores
	runtime.GOMAXPROCS(runtime.NumCPU())
	//Set default glfw controller
	controller.SetDefaultConstructor(glfwController.NewActionMap)
}

func main() {
	server := len(os.Args) > 1 && os.Args[1] == "server"
	var gameEngine engine.Engine
	var glRenderer *opengl.OpenglRenderer
	network := networking.NewNetwork()

	// Start server or connect to server
	if server {
		gameEngine = engine.NewHeadlessEngine()
		network.StartServer(serverPort)
	} else {
		glRenderer = &opengl.OpenglRenderer{
			WindowTitle:  "Networking example",
			WindowWidth:  800,
			WindowHeight: 800,
		}
		gameEngine = engine.NewEngine(glRenderer)
		network.ConnectClient(fmt.Sprintf("%v:%v", serverAddr, serverPort))
	}
	gameEngine.AddUpdatable(network)

	// map containing each player's entity
	players := make(map[string]*Player)

	//Networked Game events
	network.ClientJoinedEvent(func(clientId string) {
		fmt.Println("client joined, clientId: ", clientId)
		network.TriggerOnServerAndClients("spawn", util.SerializeArgs(clientId, startingPosition))
		for _, player := range players {
			network.TriggerEvent("spawn", clientId, util.SerializeArgs(player.clientId, player.position))
		}
	})

	network.RegisterEvent("spawn", func(clientId string, data []byte) {
		buf := bytes.NewBuffer(data)
		playerID := util.Stringfrombytes(buf)
		position := util.Vector2frombytes(buf)
		if _, ok := players[playerID]; !ok {
			player := &Player{clientId: playerID}
			players[player.clientId] = player
			player.node = renderer.NewNode()
			player.position = position
			gameEngine.AddUpdatable(player)
			if network.IsClient() {
				boxGeometry := renderer.CreateBox(30, 30)
				boxGeometry.Material = renderer.NewMaterial()
				boxGeometry.SetColor(color.NRGBA{254, 0, 0, 254})
				player.node.Add(boxGeometry)
				gameEngine.AddOrtho(player.node)
			}
		}
	})

	network.RegisterEvent("move", func(clientId string, data []byte) {
		buf := bytes.NewBuffer(data)
		playerID := util.Stringfrombytes(buf)
		velocity := util.Vector2frombytes(buf)
		if network.IsServer() && clientId != playerID {
			return // client is only allowed to control the player assigned to them.
		}
		if player, ok := players[playerID]; ok {
			player.velocity = velocity
			if network.IsServer() {
				network.BroadcastEvent("updatePlayer", util.SerializeArgs(playerID, player.position, player.velocity))
				network.FlushAllWriteBuffers()
			}
		}
	})

	network.RegisterEvent("updatePlayer", func(clientId string, data []byte) {
		if network.IsClient() { // This is a server to client update only
			buf := bytes.NewBuffer(data)
			playerID := util.Stringfrombytes(buf)
			position := util.Vector2frombytes(buf)
			velocity := util.Vector2frombytes(buf)
			if player, ok := players[playerID]; ok {
				player.position = position
				player.velocity = velocity
			}
		}
	})

	// client setup
	gameEngine.Start(func() {
		if network.IsClient() {
			gameEngine.InitFpsDial()

			if shader, err := assets.ImportShader("build/shaders/basic.vert", "build/shaders/basic.frag"); err == nil {
				gameEngine.DefaultShader(shader)
			}

			glRenderer.BackGroundColor(0, 0.4, 0, 0)

			// input/controller manager
			controllerManager := glfwController.NewControllerManager(glRenderer.Window)

			// networked movement controls
			move := func(velocity mgl32.Vec2) {
				network.TriggerOnServerAndClients("move", util.SerializeArgs(network.ClientToken(), velocity))
			}

			customController := controller.CreateController()
			controllerManager.AddController(customController.(glfwController.Controller))
			customController.BindKeyAction(func() { move(mgl32.Vec2{0, -100}) }, controller.KeyW, controller.Press)
			customController.BindKeyAction(func() { move(mgl32.Vec2{-100, 0}) }, controller.KeyA, controller.Press)
			customController.BindKeyAction(func() { move(mgl32.Vec2{0, 100}) }, controller.KeyS, controller.Press)
			customController.BindKeyAction(func() { move(mgl32.Vec2{100, 0}) }, controller.KeyD, controller.Press)
			customController.BindKeyAction(func() { move(mgl32.Vec2{}) }, controller.KeyW, controller.Release)
			customController.BindKeyAction(func() { move(mgl32.Vec2{}) }, controller.KeyA, controller.Release)
			customController.BindKeyAction(func() { move(mgl32.Vec2{}) }, controller.KeyS, controller.Release)
			customController.BindKeyAction(func() { move(mgl32.Vec2{}) }, controller.KeyD, controller.Release)
		}
	})
}
