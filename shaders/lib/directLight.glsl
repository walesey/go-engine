#frag
vec3 directLight( vec3 light, vec3 direction, vec4 diffuse, vec4 specular, vec4 normal ) {
	vec3 normal_tangentSpace = (normal.xyz*2) - 1;
	vec3 direction_tangentSpace = direction * TBNMatrix;

	float diffuseMultiplier = max(0.0, dot(normal_tangentSpace.xyz, -direction_tangentSpace));

	vec3 color = (diffuseMultiplier * diffuse.rgb);

	return color * light;
}
#endfrag
