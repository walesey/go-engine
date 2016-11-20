uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;
uniform mat4 modelNormal;
uniform vec3 cameraTranslation;

uniform bool unlit;
uniform bool useTextures;

#vert
in vec3 vert;
in vec3 normal;
in vec2 texCoord;
in vec4 color;
#endvert

#frag
out vec4 outputColor;
#endfrag

