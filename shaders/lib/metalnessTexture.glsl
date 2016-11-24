#include "./textures.glsl"

#vert
void metalnessTexture() {}
#endvert

#frag
uniform sampler2D metalnessMap;

vec4 metalness;

void metalnessTexture() {
	vec2 overflowTextCoord = repeatTextCoord();

	metalness = texture(metalnessMap, overflowTextCoord);
	specular = mix(vec4(0.04), diffuse, metalness.r);
}
#endfrag