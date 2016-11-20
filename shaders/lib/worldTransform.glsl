#include "./base.glsl"

#frag
in vec3 worldVertex;
in vec3 worldNormal;
in vec3 eyeDirection;
in mat3 TBNMatrix;
in mat3 inverseTBNMatrix;
#endfrag

#vert
out vec3 worldVertex;
out vec3 worldNormal;
out vec3 eyeDirection;
out mat3 TBNMatrix;
out mat3 inverseTBNMatrix;

void worldTransform() {
	worldVertex = (model * vec4(vert,1)).xyz;
	worldNormal = (modelNormal * vec4(normal,1)).xyz;
	eyeDirection = normalize(worldVertex - cameraTranslation);

	// generate arbitrary tangent and bitangent to the normal
	vec3 tangent = cross(normal, normal + vec3(-1));
	vec3 bitangent = cross(normal, tangent);
	vec3 worldTangent = (modelNormal * vec4(tangent,1)).xyz;
	vec3 worldBitangent = (modelNormal * vec4(bitangent,1)).xyz;

	//tangent space conversion - worldToTangent
	TBNMatrix = mat3(worldTangent, worldBitangent, worldNormal);
	inverseTBNMatrix = inverse(TBNMatrix);
}
#endvert