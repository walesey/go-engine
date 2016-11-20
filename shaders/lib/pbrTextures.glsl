#include "./textures.glsl"

#vert
void pbrTextures() {}
#endvert

#frag
uniform sampler2D roughnessMap;
uniform sampler2D metalnessMap;

vec4 roughness;
vec4 metalness;

void pbrTextures() {
	vec2 overflowTextCoord = repeatTextCoord();

	roughness = texture(roughnessMap, overflowTextCoord);
	metalness = texture(metalnessMap, overflowTextCoord);
  specular = mix(vec4(0.04), diffuse, metalness.r);
}
#endfrag