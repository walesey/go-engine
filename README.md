# GoEngine

A Simple Rendering and scene graph library for Golang.

[![Build Status](https://travis-ci.org/walesey/go-engine.svg?branch=master)](https://travis-ci.org/walesey/go-engine)

## Features
* OpenGL renderer.
* Obj importer
* Lighting Engine
* Particle System
* UI system
* Controller system (mouse, keyboard, joystick)
* Map editor
* multiplayer networking library

## Instructions
Example programs can be found in examples/*
uses go-gl which requires opengl.

## Core Packages
* renderer - package contains common renderer interface and scenegraph implementation.
* opengl - package contains opengl renderer implementation.
* engine - package Is the high level engine interface that handles a lot of boilerplate stuff.
* controller - package Is the api for keyboard/mouse/joystick controllers. (see examples/simple/main.go)
* assets - asset management for images and obj files.
* vectormath - package contains useful vector math functions and types used throughout the engine.

## Important Interfaces and Structs
* renderer.Entity (interface) - anything that can be moved, rotated and scaled. (eg. Camera/Node/ParticleEmitter)
* renderer.Spatial (interface) - something that can be Drawn by a Renderer (eg. Geometry/Node)
* renderer.Node (struct) - Container for Spatials.
* renderer.Geometry (struct) - A collection of faces and verticies.
* renderer.Material (struct) - used for texturing a geometry.
* renderer.Camera (struct) - The main object used to interact with the camera.
* controller.Controller (interface) - Can have (mouse/keyboard...) events bound to.
* vectormath.Vector3 (struct) - used to describe position in 3D space (X,Y,Z)
* vectormath.Quaternion (struct) - used to describe orientation in 3D space (X,Y,Z,W)
* engine.Engine (interface) - The main game engine interface
* engine.Updatable (interface) - anything that can be updated every game simulation step.
