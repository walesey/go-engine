# GoEngine

Go engine is a simple game engine intended for 2D or 3D games.

## Features

- OpenGL renderer
- Obj importer
- Lighting Engine
- Particle System
- UI system
- Controller system (mouse, keyboard, joystick)
- Multiplayer networking library

## Platform Support

- Windows
- macOS

## Instructions

Example programs can be found in `examples/*`.

Installing deps:
sudo apt-get install mesa-utils
sudo apt install mesa-common-dev
sudo apt-get install libx11-dev
sudo apt-get install libglfw3
sudo apt-get install libxrandr-dev libxcursor-dev libxinerama-dev libxi-dev
sudo apt-get install libxxf86vm-dev

## Core Packages

- renderer - package contains common renderer interface and scenegraph implementation.
- opengl - package contains opengl renderer implementation.
- engine - package Is the high level engine interface that handles a lot of boilerplate stuff.
- controller - package Is the api for keyboard/mouse/joystick controllers. (see examples/simple/main.go)
- assets - asset management for images and obj files.

## Important Interfaces and Structs

- renderer.Entity (interface) - anything that can be moved, rotated and scaled. (eg. Camera/Node/ParticleEmitter)
- renderer.Spatial (interface) - something that can be Drawn by a Renderer (eg. Geometry/Node)
- renderer.Node (struct) - Container for Spatials.
- renderer.Geometry (struct) - A collection of faces and verticies.
- renderer.Material (struct) - used for texturing a geometry.
- renderer.Camera (struct) - Struct used to manage the camera.
- renderer.Light (struct) - Struct used to manage dynamic lights.
- controller.Controller (interface) - Can have (mouse/keyboard...) events bound to.
- engine.Engine (interface) - The main game engine interface
- engine.Updatable (interface) - anything that can be updated every game simulation step.

![Demo](http://i.imgur.com/toTtrxp.jpg)
