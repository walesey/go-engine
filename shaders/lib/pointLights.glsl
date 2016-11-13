#include "./base.glsl"
#include "./worldTransform.glsl"
#include "./directLight.glsl"

#define MAX_LIGHTS 4

uniform int nbPointLights;
uniform vec4 pointLightPositions[ MAX_LIGHTS ];
uniform vec4 pointLightValues[ MAX_LIGHTS ];

#frag
vec4 pointLights(vec4 diffuse, vec4 specular, vec4 normal) {
	vec4 totalLight = vec4(0.1, 0.0, 0.0, 1.0);
	for (int i=0; i < nbPointLights; i++) {
		vec4 LightPos = pointLightPositions[i];
		vec4 LightValue = pointLightValues[i];

		vec4 v = vec4( worldVertex - LightPos );
		float lightDistance = v.x*v.x + v.y*v.y + v.z*v.z;
		float brightness = 1.0 / lightDistance;
		vec4 worldLightDir = normalize(v);

		totalLight += directLight(brightness, worldLightDir, diffuse, specular, normal);
	}
	return totalLight;
}
#endfrag