#version 330

#include "./lib/base.glsl"
#include "./lib/worldTransform.glsl"
#include "./lib/pointLights.glsl"

uniform sampler2D normal;
uniform sampler2D diffuse;
uniform sampler2D ao;
uniform sampler2D metalness;
uniform sampler2D roughness;

void main() {
	#vert
	worldTransform();
	gl_Position = projection * camera * vec4(worldVertex, 1);
	#endvert

	#frag
	vec4 finalColor = vec4();

	finalColor += pointLightColor();
	
	output = finalColor;
	#endfrag
}