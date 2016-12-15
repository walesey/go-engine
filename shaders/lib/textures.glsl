#include "./base.glsl"

#vert
out vec2 fragTexCoord;
out vec4 fragColor;

void textures() {
	fragTexCoord = texCoord;
	fragColor = color;
}
#endvert

#frag
uniform sampler2D normalMap;
uniform sampler2D diffuseMap;
uniform sampler2D specularMap;
uniform sampler2D aoMap;

in vec2 fragTexCoord;
in vec4 fragColor;

vec4 normalValue;
vec4 diffuse;
vec4 specular;
vec4 ao;

vec2 repeatTextCoord() {
	float textureX = fragTexCoord.x - int(fragTexCoord.x);
	float textureY = fragTexCoord.y - int(fragTexCoord.y);
	if (fragTexCoord.x < 0) {textureX = textureX + 1.0;}
	if (fragTexCoord.y < 0) {textureY = textureY + 1.0;}
	return vec2(textureX, textureY);
}

void textures() {
	vec2 overflowTextCoord = repeatTextCoord();
	
	// multiply color by diffuse map. use only color if no map is provided
	if (useTextures) {
		diffuse = fragColor * texture(diffuseMap, overflowTextCoord);
		specular = texture(specularMap, overflowTextCoord);
		normalValue = texture(normalMap, overflowTextCoord);
		ao = texture(aoMap, overflowTextCoord);
	} else {
		diffuse = fragColor;
		specular = vec4(0);
		normalValue = vec4(0);
		ao = vec4(1);
	}
}
#endfrag