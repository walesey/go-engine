#include "./textures.glsl"

#vert
void pbrCompositeTextures() {}
#endvert

#frag
uniform sampler2D compositeMap;

vec4 roughness;
vec4 metalness;
vec4 metalSpecular;
vec4 metalDiffuse;

void pbrCompositeTextures() {
	vec2 overflowTextCoord = repeatTextCoord();

	vec4 composite = texture(compositeMap, overflowTextCoord);

	ao = vec4(composite.r);
	roughness = vec4(composite.g);
	metalness = vec4(composite.b);

	metalSpecular = mix(vec4(0.04), diffuse, metalness.r);
	metalDiffuse = mix(diffuse, vec4(0), metalness.r);
}
#endfrag