#version 330

uniform sampler2D diffuseMap;
in vec2 fragTexCoord;
out vec4 outputColor;

void main() {

	vec4 finalColor = vec4(0,0,0,1);
	vec4 source = texture(diffuseMap, fragTexCoord);
	finalColor = source;
	outputColor = vec4(finalColor.xyz, 1.0); 
  
}
