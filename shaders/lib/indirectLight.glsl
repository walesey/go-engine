#frag
#include "./worldTransform.glsl"

uniform samplerCube environmentMap;

vec3 indirectLight(vec4 diffuse, vec4 specular, vec4 normal) {
	vec3 normal_tangentSpace = (normal.xyz*2) - 1;
	vec3 normal_worldSpace = normal_tangentSpace * inverseTBNMatrix;
  vec3 diffuseValue = texture(environmentMap, normal_worldSpace).xyz;

	return (10*diffuseValue); //+ (specular * specularValue);
}
#endfrag