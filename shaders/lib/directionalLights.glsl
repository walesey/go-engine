#frag
#include "./base.glsl"
#include "./worldTransform.glsl"
#include "./directLight.glsl"

#define MAX_DIRECTIONAL_LIGHTS 4

uniform int nbDirectionalLights;
uniform vec4 directionalLightVectors[ MAX_DIRECTIONAL_LIGHTS ];
uniform vec4 directionalLightValues[ MAX_DIRECTIONAL_LIGHTS ];

vec3 directionalLights(vec4 diffuse, vec4 specular, vec4 normal) {
	vec3 totalLight = vec3(0.0, 0.0, 0.0);
	for (int i=0; i < nbDirectionalLights; i++) {
		vec3 LightDirection = directionalLightVectors[i].rgb;
		vec3 LightValue = directionalLightValues[i].rgb;

		totalLight += directLight(LightValue, LightDirection, diffuse, specular, normal);
	}
	return totalLight;
}
#endfrag