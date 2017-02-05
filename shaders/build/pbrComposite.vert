#version 400

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;
uniform mat4 modelNormal;
uniform vec3 cameraTranslation;

uniform bool unlit;
uniform bool useTextures;

in vec3 vert;
in vec3 normal;
in vec2 texCoord;
in vec4 color;

out vec3 worldVertex;
out vec3 worldNormal;
out vec3 eyeDirection;
out mat3 TBNMatrix;
out mat3 inverseTBNMatrix;

void worldTransform() {
	worldVertex = (model * vec4(vert,1)).xyz;
	worldNormal = (modelNormal * vec4(normal,1)).xyz;
	worldNormal = normalize(worldNormal);
	eyeDirection = normalize(worldVertex - cameraTranslation);

	// generate arbitrary tangent and bitangent to the normal
	vec3 tangent = cross(normal, normal + vec3(-1));
	vec3 bitangent = cross(normal, tangent);
	vec3 worldTangent = normalize((modelNormal * vec4(tangent,1)).xyz);
	vec3 worldBitangent = normalize((modelNormal * vec4(bitangent,1)).xyz);

	//tangent space conversion - worldToTangent
	TBNMatrix = mat3(worldTangent, worldBitangent, worldNormal);
	inverseTBNMatrix = inverse(TBNMatrix);
}

out vec2 fragTexCoord;
out vec4 fragColor;

void textures() {
	fragTexCoord = texCoord;
	fragColor = color;
}

void pbrCompositeTextures() {}

float pow2(float x) { 
	return x*x; 
}

float pow3(float x) { 
	return x*x*x; 
}

void glowOutput() {}

void main() {
  textures();
	pbrCompositeTextures();
	glowOutput();

	worldTransform();
	gl_Position = projection * camera * model * vec4(vert, 1);

}
