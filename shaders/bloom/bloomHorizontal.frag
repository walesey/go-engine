#version 330

//material
uniform sampler2D diffuse;

uniform vec2 size;
uniform float threshold;
uniform float intensity;
uniform float quality;
uniform int samples;

in vec2 fragTexCoord;

out vec4 outputColor;

void main() {
	vec4 finalColor = vec4(0,0,0,1);

	int diff = (samples - 1) / 2;
	vec2 sizeFactor = vec2(1,1) / size * quality;
	vec4 source = texture(diffuse, fragTexCoord);

	for (int x = -diff; x <= diff; x++){
		vec2 offset = vec2(x, 0) * sizeFactor;
		vec4 fragment = texture(diffuse, fragTexCoord + offset);
		if (fragment.r > threshold || fragment.g > threshold || fragment.b > threshold) {
			finalColor += fragment;
		}
	}

	finalColor = intensity*(finalColor / samples) + source;

	//final output
	outputColor = finalColor; 
}
