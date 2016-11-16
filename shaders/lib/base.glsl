uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;
uniform mat4 modelNormal;

uniform bool unlit;
uniform bool useTextures;

#vert
in vec3 vert;
in vec3 normal;
in vec2 texCoord;
in vec4 color;

out vec2 fragTexCoord;
out vec4 fragColor;
#endvert

#frag
in vec2 fragTexCoord;
in vec4 fragColor;

out vec4 outputColor;
#endfrag

