#include "./base.glsl"

#vert
out vec2 fragTexCoord;
#endvert

#frag
in vec2 fragTexCoord;

uniform sampler2D normalMap;
uniform sampler2D diffuseMap;
uniform sampler2D specularMap;
uniform sampler2D aoMap;

vec4 normal;
vec4 diffuse;
vec4 specular;
vec4 ao;
#endfrag

void textures() {
	#vert
	fragTexCoord = texCoord;
	#endvert
	
	#frag
	normal = texture(normalMap, fragTexCoord);
	diffuse = texture(diffuseMap, fragTexCoord);
	specular = texture(specularMap, fragTexCoord);
	ao = texture(aoMap, fragTexCoord);
	#endfrag
}