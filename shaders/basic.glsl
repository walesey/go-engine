#version 400

#include "./lib/base.glsl"
#include "./lib/worldTransform.glsl"
#include "./lib/textures.glsl"
#include "./lib/ambientLight.glsl"
#include "./lib/pointLights.glsl"
#include "./lib/directionalLights.glsl"

void main() {
	textures();

	#vert
	worldTransform();
	gl_Position = projection * camera * model * vec4(vert, 1);
	#endvert

	#frag
	if (unlit) {
		outputColor = diffuse;
	} else {
		vec3 finalColor = ambientLight(diffuse) + pointLights(diffuse, specular, normal) + directionalLights(diffuse, specular, normal);
		outputColor = vec4(finalColor, diffuse.a);
	}
	#endfrag
}