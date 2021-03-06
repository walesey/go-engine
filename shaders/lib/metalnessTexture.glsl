#include "./textures.glsl"

#vert
void metalnessTexture() {}
#endvert

#frag
uniform sampler2D metalnessMap;

vec4 metalness;
vec4 metalSpecular;
vec4 metalDiffuse;

void metalnessTexture() {
	vec2 overflowTextCoord = repeatTextCoord();

	metalness = texture(metalnessMap, overflowTextCoord);
	metalSpecular = mix(vec4(0.04), diffuse, metalness.r);
	metalDiffuse = mix(diffuse, vec4(0), metalness.r);
}
#endfrag