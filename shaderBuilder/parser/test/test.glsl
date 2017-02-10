#version 330

#include "includeTest.glsl"
#include "includeTest2.glsl"

void main() {
#vert
  gl_Position = vertFunc();
#endvert

#frag
  gl_FragColor = fragFunc();
#endfrag
}
