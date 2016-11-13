#frag
vec4 directLight( float brightness, vec4 direction, vec4 diffuse, vec4 specular, vec4 normal ) {
	vec4 normal_tangentSpace = (normal*2) - 1;
	vec3 direction_tangentSpace = direction.xyz * TBNMatrix;

	float diffuseMultiplier = max(0.0, dot(normal_tangentSpace.xyz, -direction_tangentSpace.xyz));

	return (diffuseMultiplier * diffuse);
}
#endfrag