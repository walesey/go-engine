#version 330

uniform sampler2D tex0;
in vec2 fragTexCoord;
out vec4 outputColor;

float cell( float source ) {
	if (source < 0.3) {
		return 0;
	} else if (source < 0.6) {
		return 0.3;
	} else if (source < 0.9) {
		return 0.6;
	} 
	return 1;
}

void main() {

	vec4 finalColor = vec4(0,0,0,1);
	vec4 source = texture(tex0, fragTexCoord);
	finalColor = vec4( cell(source.r), cell(source.g), cell(source.b), cell(source.a) );
	outputColor = finalColor; 
  
}
