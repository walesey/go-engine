#version 330

//material
uniform sampler2D diffuse;

uniform vec2 size;
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
		finalColor += texture(diffuse, fragTexCoord + offset);
	}

	finalColor = (0.5*(finalColor / samples) + 0.5*source);

	//final output
	outputColor = finalColor; 
}
