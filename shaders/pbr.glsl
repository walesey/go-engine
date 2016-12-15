#version 400

#include "./lib/base.glsl"
#include "./lib/worldTransform.glsl"
#include "./lib/textures.glsl"
#include "./lib/metalnessTexture.glsl"
#include "./lib/roughnessTexture.glsl"
#include "./lib/fresnelEffect.glsl"
#include "./lib/ambientLight.glsl"
#include "./lib/pointLights.glsl"
#include "./lib/directionalLights.glsl"
#include "./lib/indirectLight.glsl"

void main() {
	textures();
	metalnessTexture();
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
		vec4 feSpecular = fresnelEffect(specular, normalValue);
		vec3 dLight = ambientLight(aoDiffuse) + pointLights(aoDiffuse, feSpecular, normalValue) + directionalLights(aoDiffuse, feSpecular, normalValue);
		vec3 iLight = indirectLight(aoDiffuse, feSpecular, normalValue);
		outputColor = vec4(dLight + iLight, diffuse.a);
	}
	#endfrag
}