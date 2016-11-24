#include "./textures.glsl"

#vert
void roughnessTexture() {}
#endvert

#frag
uniform sampler2D roughnessMap;

vec4 roughness;

void roughnessTexture() {
	vec2 overflowTextCoord = repeatTextCoord();

	roughness = texture(roughnessMap, overflowTextCoord);
}
#endfrag