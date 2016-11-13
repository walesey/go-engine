#version 330

#include "./lib/base.glsl"
#include "./lib/worldTransform.glsl"
#include "./lib/textures.glsl"
#include "./lib/pointLights.glsl"

void main() {
	textures();

	#vert
	worldTransform();
	gl_Position = projection * camera * model * vec4(vert, 1);
	#endvert

	#frag
	vec4 finalColor = vec4(0.0);
	finalColor += ao * pointLights(diffuse, specular, normal);
	outputColor = finalColor;
	#endfrag
}