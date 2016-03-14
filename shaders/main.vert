#version 330

uniform vec4 worldCamPos;
uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;

in vec3 vert;
in vec3 normal;
in vec3 tangent;
in vec3 bitangent;
in vec2 vertTexCoord;
in vec4 color;

out mat3 TBNMatrix;
out mat3 inverseTBNMatrix;
out vec4 worldVertex;
out vec3 worldNormal;
out vec2 fragTexCoord;
out vec4 fragColor;
out vec3 tangentEyeDirection;

void main() {
	worldVertex = model * vec4(vert, 1);
	gl_Position = projection * camera * worldVertex;
	mat4 modelNormal = transpose(inverse(model));
	worldNormal = (modelNormal * vec4(normal,1)).xyz;
	vec3 worldTangent = (modelNormal * vec4(tangent,1)).xyz;
	vec3 worldBitangent = (modelNormal * vec4(bitangent,1)).xyz;
	//tangent space conversion - worldToTangent
	TBNMatrix = mat3(worldTangent, worldBitangent, worldNormal);
	inverseTBNMatrix = inverse(TBNMatrix);
	vec4 worldEyeDirection = vec4( worldVertex - worldCamPos );
	tangentEyeDirection = normalize( worldEyeDirection.xyz * TBNMatrix );
	fragTexCoord = vertTexCoord;
	fragColor = color;
}
