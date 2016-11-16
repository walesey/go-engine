#include "./base.glsl"

#vert
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

vec4 normal;
vec4 diffuse;
vec4 specular;
vec4 ao;

void textures() {
	// repeat textures
	float textureX = fragTexCoord.x - int(fragTexCoord.x);
	float textureY = fragTexCoord.y - int(fragTexCoord.y);
	if (fragTexCoord.x < 0) {textureX = textureX + 1.0;}
	if (fragTexCoord.y < 0) {textureY = textureY + 1.0;}
	vec2 overflowTextCoord = vec2(textureX, textureY);
	
	// multiply color by diffuse map. use only color if no map is provided
	if (useTextures) {
		diffuse = fragColor * texture(diffuseMap, overflowTextCoord);
	} else {
		diffuse = fragColor;
	}

	normal = texture(normalMap, overflowTextCoord);
	specular = texture(specularMap, overflowTextCoord);
	ao = texture(aoMap, overflowTextCoord);
}
#endfrag