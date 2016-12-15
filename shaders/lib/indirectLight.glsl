#frag
#include "./worldTransform.glsl"

uniform samplerCube environmentMap;

vec3 indirectLight(vec4 diffuse, vec4 specular, vec4 normalValue) {
	vec3 normal_tangentSpace = (normalValue.xyz*2) - 1;
	vec3 normal_worldSpace = normal_tangentSpace * inverseTBNMatrix;
	vec3 reflectedEye_worldSpace = reflect( eyeDirection, normal_worldSpace );

  vec3 diffuseValue = textureLod(environmentMap, normal_worldSpace, 10).rgb;
  vec3 specularValue = textureLod(environmentMap, reflectedEye_worldSpace, roughness.r * 10).rgb;

	return (diffuse.rgb * diffuseValue) + (specular.rgb * specularValue);
}
#endfrag