#frag
#include "./common.glsl"
#include "./worldTransform.glsl"

vec4 fresnelEffect(vec4 baseSpecular, vec4 normalValue) {
	vec3 normal_tangentSpace = (normalValue.xyz*2) - 1;
	vec3 normal_worldSpace = normal_tangentSpace * inverseTBNMatrix;
	float NdV = abs(dot(normal_worldSpace, eyeDirection));

  return mix(baseSpecular, vec4(1.0), pow(1.0 - NdV, 5.0));
}
#endfrag