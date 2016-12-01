#version 400

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;
uniform mat4 modelNormal;
uniform vec3 cameraTranslation;

uniform bool unlit;
uniform bool useTextures;

out vec4 outputColor;

in vec3 worldVertex;
in vec3 worldNormal;
in vec3 eyeDirection;
in mat3 TBNMatrix;
in mat3 inverseTBNMatrix;

uniform sampler2D normalMap;
uniform sampler2D diffuseMap;
uniform sampler2D specularMap;
uniform sampler2D aoMap;

in vec2 fragTexCoord;
in vec4 fragColor;

vec4 normal;
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
		normal = texture(normalMap, overflowTextCoord);
		ao = texture(aoMap, overflowTextCoord);
	} else {
		diffuse = fragColor;
		specular = vec4(0);
		normal = vec4(0);
		ao = vec4(1);
	}
}

float pow2(float x) { 
	return x*x; 
}

float pow3(float x) { 
	return x*x*x; 
}

vec4 fresnelEffect(vec4 baseSpecular, vec4 normal) {
	vec3 normal_tangentSpace = (normal.xyz*2) - 1;
	vec3 normal_worldSpace = normal_tangentSpace * inverseTBNMatrix;
	float NdV = abs(dot(normal_worldSpace, eyeDirection));

  return mix(baseSpecular, vec4(1.0), pow(1.0 - NdV, 5.0));
}

uniform sampler2D roughnessMap;

vec4 roughness;

void roughnessTexture() {
	vec2 overflowTextCoord = repeatTextCoord();

	roughness = texture(roughnessMap, overflowTextCoord);
}

vec3 directLight( vec3 light, vec3 direction, vec4 diffuse, vec4 specular, vec4 normal ) {
	vec3 normal_tangentSpace = (normal.xyz*2) - 1;
	vec3 direction_tangentSpace = direction * TBNMatrix;
	vec3 eyeDirection_tangentSpace = eyeDirection * TBNMatrix;
	vec3 reflectedEye_tangentSpace = reflect( eyeDirection_tangentSpace, normal_tangentSpace );

	float diffuseMultiplier = max(0.0, dot(normal_tangentSpace, -direction_tangentSpace));

	float specularMultiplier = pow2(max(0.0, dot(reflectedEye_tangentSpace, -direction_tangentSpace)));

	vec3 color = (diffuseMultiplier * diffuse.rgb) + (specularMultiplier * specular.rgb);

	return color * light;
}

#define MAX_POINT_LIGHTS 4

uniform int nbPointLights;
uniform vec4 pointLightPositions[ MAX_POINT_LIGHTS ];
uniform vec4 pointLightValues[ MAX_POINT_LIGHTS ];

vec3 pointLights(vec4 diffuse, vec4 specular, vec4 normal) {
	vec3 totalLight = vec3(0.0, 0.0, 0.0);
	for (int i=0; i < nbPointLights; i++) {
		vec3 LightPos = pointLightPositions[i].rgb;
		vec3 LightValue = pointLightValues[i].rgb;

		vec3 v = worldVertex - LightPos;
		float lightDistance = dot(v, v);
		float brightness = 1.0 / lightDistance;

		vec3 worldLightDir = normalize(v);
		vec3 light = brightness*LightValue;

		totalLight += directLight(light, worldLightDir, diffuse, specular, normal);
	}
	return totalLight;
}

#define MAX_DIRECTIONAL_LIGHTS 4

uniform int nbDirectionalLights;
uniform vec4 directionalLightVectors[ MAX_DIRECTIONAL_LIGHTS ];
uniform vec4 directionalLightValues[ MAX_DIRECTIONAL_LIGHTS ];

vec3 directionalLights(vec4 diffuse, vec4 specular, vec4 normal) {
	vec3 totalLight = vec3(0.0, 0.0, 0.0);
	for (int i=0; i < nbDirectionalLights; i++) {
		vec3 LightDirection = directionalLightVectors[i].rgb;
		vec3 LightValue = directionalLightValues[i].rgb;

		totalLight += directLight(LightValue, LightDirection, diffuse, specular, normal);
	}
	return totalLight;
}

uniform samplerCube environmentMap;

vec3 indirectLight(vec4 diffuse, vec4 specular, vec4 normal) {
	vec3 normal_tangentSpace = (normal.xyz*2) - 1;
	vec3 normal_worldSpace = normal_tangentSpace * inverseTBNMatrix;
	vec3 reflectedEye_worldSpace = reflect( eyeDirection, normal_worldSpace );

  vec3 diffuseValue = textureLod(environmentMap, normal_worldSpace, 10).rgb;
  vec3 specularValue = textureLod(environmentMap, reflectedEye_worldSpace, roughness.r * 10).rgb;

	return (diffuse.rgb * diffuseValue) + (specular.rgb * specularValue);
}

void main() {
	textures();
	roughnessTexture();

	if (unlit) {
		outputColor = diffuse;
	} else {
		vec4 aoDiffuse = ao * diffuse;
		vec4 feSpecular = fresnelEffect(specular, normal);
		vec3 dLight = pointLights(aoDiffuse, feSpecular, normal) + directionalLights(aoDiffuse, feSpecular, normal);
		vec3 iLight = indirectLight(aoDiffuse, feSpecular, normal);
		outputColor = vec4(dLight + iLight, diffuse.a);
	}
	
}
