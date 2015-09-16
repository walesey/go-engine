#version 330

//material
uniform sampler2D diffuse;

in vec2 fragTexCoord;

out vec4 outputColor;

void main() {
	vec4 finalColor = vec4(0,0,0,1);

	//TODO: move these to uniforms
	vec2 size = vec2(1900, 1000);
	float quality = 3.5; //lower = higher quality and less blur
	int samples = 25;

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
