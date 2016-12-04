#frag
#include "./base.glsl"

uniform vec3 ambientLightValue;

vec3 ambientLight(vec4 diffuse) {
	return ambientLightValue * diffuse.rgb;
}
#endfrag