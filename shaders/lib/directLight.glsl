#frag
#include "./common.glsl"
#include "./worldTransform.glsl"

vec3 directLight( vec3 light, vec3 direction, vec4 diffuse, vec4 specular, vec4 normalValue ) {
	vec3 normal_tangentSpace = (normalValue.xyz*2) - 1;
	vec3 direction_tangentSpace = direction * TBNMatrix;
	vec3 eyeDirection_tangentSpace = eyeDirection * TBNMatrix;
	vec3 reflectedEye_tangentSpace = reflect( eyeDirection_tangentSpace, normal_tangentSpace );

	float diffuseMultiplier = max(0.0, dot(normal_tangentSpace, -direction_tangentSpace));

	float specularMultiplier = pow2(max(0.0, dot(reflectedEye_tangentSpace, -direction_tangentSpace)));

	vec3 color = (diffuseMultiplier * diffuse.rgb) + (specularMultiplier * specular.rgb);

	return color * light;
}
#endfrag
