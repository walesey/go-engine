#frag
#include "./base.glsl"
#include "./worldTransform.glsl"
#include "./directLight.glsl"

#define MAX_POINT_LIGHTS 4

uniform int nbPointLights;
uniform vec4 pointLightPositions[ MAX_POINT_LIGHTS ];
uniform vec4 pointLightValues[ MAX_POINT_LIGHTS ];

vec3 pointLights(vec4 diffuse, vec4 specular, vec4 normalValue) {
	vec3 totalLight = vec3(0.0, 0.0, 0.0);
	for (int i=0; i < nbPointLights; i++) {
		vec3 LightPos = pointLightPositions[i].rgb;
		vec3 LightValue = pointLightValues[i].rgb;

		vec3 v = worldVertex - LightPos;
		float lightDistance = dot(v, v);
		float brightness = 1.0 / lightDistance;

		vec3 worldLightDir = normalize(v);
		vec3 light = brightness*LightValue;

		totalLight += directLight(light, worldLightDir, diffuse, specular, normalValue);
	}
	return totalLight;
}
#endfrag