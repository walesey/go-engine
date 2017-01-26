#version 330

#vert
in vec2 vert;
in vec2 texCoord;
out vec2 fragTexCoord;
#endvert

#frag
uniform sampler2D diffuseMap;
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
#endfrag

void main() {
	#vert
	gl_Position = vec4(vert, 0.0, 1.0);
	fragTexCoord = texCoord;
	#endvert

  #frag
	vec4 finalColor = vec4(0,0,0,1);
	vec4 source = texture(diffuseMap, fragTexCoord);
	finalColor = vec4( cell(source.r), cell(source.g), cell(source.b), cell(source.a) );
	outputColor = finalColor; 
  #endfrag
}