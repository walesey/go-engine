#version 330

#include "includeTest.glsl"

void main() {
  #vert
  gl_Position = vertFunc();
  #endvert

  #frag
  gl_FragColor = fragFunc();
  #endfrag
}