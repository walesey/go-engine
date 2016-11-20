#version 330

#include "./lib/base.glsl"
#include "./lib/worldTransform.glsl"
#include "./lib/textures.glsl"
#include "./lib/pbrTextures.glsl"
#include "./lib/pointLights.glsl"
#include "./lib/directionalLights.glsl"
#include "./lib/indirectLight.glsl"

void main() {
	textures();
	pbrTextures();

	#vert
	worldTransform();
	gl_Position = projection * camera * model * vec4(vert, 1);
	#endvert

	#frag
	if (unlit) {
		outputColor = diffuse;
	} else {
		vec3 finalColor = vec3(0.0);
		finalColor += ao.rgb * pointLights(diffuse, specular, normal);
		finalColor += ao.rgb * directionalLights(diffuse, specular, normal);
		finalColor += ao.rgb * indirectLight(diffuse, specular, normal);
		outputColor = vec4(finalColor, diffuse.a);
	}
	#endfrag
}