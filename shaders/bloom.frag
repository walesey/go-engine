#version 330

//material
uniform sampler2D diffuse;

in vec2 fragTexCoord;

out vec4 outputColor;

void main() {
	vec4 finalColor = vec4(0,0,0,1);
	
	finalColor = texture(diffuse, fragTexCoord);

	//final output
	outputColor = finalColor; 
}
