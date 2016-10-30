# GoEngine

Go engine is a simple game engine intended for 2D or 3D games.
This project is currently in progress and it is likely to change significantly over the next year or so.

[![Build Status](https://travis-ci.org/walesey/go-engine.svg?branch=master)](https://travis-ci.org/walesey/go-engine)


## Features
* OpenGL renderer
* Obj importer
* Lighting Engine
* Particle System
* UI system
* Controller system (mouse, keyboard, joystick)
* Multiplayer networking library (event driven)


## Instructions

Example programs can be found in `examples/*`.

Go-Engine uses go-gl which requires opengl.

## Core Packages
* renderer - package contains common renderer interface and scenegraph implementation.
* opengl - package contains opengl renderer implementation.
* engine - package Is the high level engine interface that handles a lot of boilerplate stuff.
* controller - package Is the api for keyboard/mouse/joystick controllers. (see examples/simple/main.go)
* assets - asset management for images and obj files.

## Important Interfaces and Structs
* renderer.Entity (interface) - anything that can be moved, rotated and scaled. (eg. Camera/Node/ParticleEmitter)
* renderer.Spatial (interface) - something that can be Drawn by a Renderer (eg. Geometry/Node)
* renderer.Node (struct) - Container for Spatials.
* renderer.Geometry (struct) - A collection of faces and verticies.
* renderer.Material (struct) - used for texturing a geometry.
* renderer.Camera (struct) - The main object used to interact with the camera.
* controller.Controller (interface) - Can have (mouse/keyboard...) events bound to.
* engine.Engine (interface) - The main game engine interface
* engine.Updatable (interface) - anything that can be updated every game simulation step.
