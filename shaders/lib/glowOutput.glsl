#include "./base.glsl"
#include "./textures.glsl"

#frag
uniform sampler2D glowMap;

out vec4 brightColor;
#endfrag

#vert
void glowOutput() {}
#endvert

#frag
void glowOutput() {
	vec2 overflowTextCoord = repeatTextCoord();

	brightColor = fragColor * texture(glowMap, overflowTextCoord);
}
#endfrag