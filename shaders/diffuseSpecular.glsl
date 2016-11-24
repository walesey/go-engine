#version 400

#include "./lib/base.glsl"
#include "./lib/worldTransform.glsl"
#include "./lib/textures.glsl"
#include "./lib/roughnessTexture.glsl"
#include "./lib/pointLights.glsl"
#include "./lib/directionalLights.glsl"
#include "./lib/indirectLight.glsl"

void main() {
	textures();
	roughnessTexture();

	#vert
	worldTransform();
	gl_Position = projection * camera * model * vec4(vert, 1);
	#endvert

	#frag
	if (unlit) {
		outputColor = diffuse;
	} else {
		vec4 aoDiffuse = ao * diffuse;
		vec3 dLight = pointLights(aoDiffuse, specular, normal) + directionalLights(aoDiffuse, specular, normal);
		vec3 iLight = indirectLight(aoDiffuse, specular, normal);
		outputColor = vec4(dLight + iLight, diffuse.a);
	}
	#endfrag
}