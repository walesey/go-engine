package shaders

var MainFrag = `
#version 330

#define MAX_LIGHTS 8

#define LIGHT_POSITION 0
#define LIGHT_AMBIENT 1
#define LIGHT_DIFFUSE 2
#define LIGHT_SPECULAR 3

#define MODE_UNLIT 0
#define MODE_LIT 1
#define MODE_EMIT 2

uniform int mode;
uniform bool useVertexColor;

uniform int nbLights;
uniform int nbDirectionalLights;
uniform vec4 lights[ MAX_LIGHTS * 4 ];
uniform vec4 directionalLights[ MAX_LIGHTS * 4 ];

//material
uniform sampler2D diffuse;
uniform sampler2D normal;
uniform sampler2D roughness;
uniform sampler2D metalness;
uniform samplerCube environmentMap;
uniform samplerCube environmentMapLOD1;
uniform samplerCube environmentMapLOD2;
uniform samplerCube environmentMapLOD3;
uniform samplerCube illuminanceMap;

uniform vec4 worldCamPos;
uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;

in mat3 TBNMatrix;
in mat3 inverseTBNMatrix;
in vec4 worldVertex;
in vec3 worldNormal;
in vec2 fragTexCoord;
in vec4 fragColor;
in vec3 tangentEyeDirection;

out vec4 outputColor;

vec4 vectorCap( vec4 vector, float cap ){
	if (vector.r > cap) {vector.r = cap;}
	if (vector.g > cap) {vector.g = cap;}
	if (vector.b > cap) {vector.b = cap;}
	return vector;
}

vec4 directBRDF( vec4 LightAmb, vec4 LightDiff, vec4 LightSpec, vec4 LightDir, vec4 diffuseValue, vec4 specularValue, vec4 tangentNormal, vec4 tangentReflectedEye){
	vec3 tangentLightDirection = LightDir.xyz * TBNMatrix;
	tangentLightDirection = normalize( tangentLightDirection );

	vec4 ambientOut = diffuseValue * LightAmb;

	float diffuseMultiplier = max(0.0, dot(tangentNormal.xyz, -tangentLightDirection));
	vec4 diffuseOut = diffuseValue * diffuseMultiplier * LightDiff;

	float specularMultiplier = pow( max(0.0, dot( tangentReflectedEye.xyz, -tangentLightDirection)), 2.0);
	vec4 specularOut = vec4( specularValue.rgb, 1 ) * specularMultiplier * LightSpec;

	return ambientOut + diffuseOut + specularOut;
}

void main() {
	//repeat textures
	float textureX = fragTexCoord.x - int(fragTexCoord.x);
	float textureY = fragTexCoord.y - int(fragTexCoord.y);
	if (fragTexCoord.x < 0) {textureX = textureX + 1.0;}
	if (fragTexCoord.y < 0) {textureY = textureY + 1.0;}
	vec2 textCoord = vec2(textureX, textureY);

	vec4 albedoValue = texture(diffuse, textCoord) * fragColor;
	if (useVertexColor) {
		albedoValue = fragColor;
	}
	vec4 normalValue = texture(normal, textCoord);
	float roughnessValue = texture(roughness, textCoord).r;
	float metalnessValue = texture(metalness, textCoord).r;
	float alphaValue = albedoValue.a;

	//metalness - if metal use albedo value as spec value 
	vec4 specularValue = mix( vec4(0.04), albedoValue, metalnessValue);
	vec4 diffuseValue = (1.0-metalnessValue) * albedoValue;

	vec4 finalColor = vec4(0,0,0,1);
	vec4 directColor = vec4(0,0,0,1);
	vec4 indirectColor = vec4(0,0,0,1);

 	if( mode == MODE_LIT ){

		//Normal calculations
		vec4 tangentNormal = normalValue * 2 - 1;
		if( abs(tangentNormal.x) < 0.1 && abs(tangentNormal.y) < 0.1 && abs(tangentNormal.z) < 0.1 ){
			tangentNormal = vec4(0,0,1,1);
		}

		//reflected eye
		vec4 tangentReflectedEye = vec4( reflect( tangentEyeDirection, tangentNormal.xyz ), 1);
		vec4 worldReflectedEye = vec4( tangentReflectedEye.xyz * inverseTBNMatrix , 1);

		//point lights
		for (int i=0;i<nbLights;i++){

			//light components
			vec4 LightPos = lights[(i*4)+LIGHT_POSITION];
			vec4 LightAmb = lights[(i*4)+LIGHT_AMBIENT];
			vec4 LightDiff = lights[(i*4)+LIGHT_DIFFUSE];
			vec4 LightSpec = lights[(i*4)+LIGHT_SPECULAR];

			//point light source
			vec4 worldLightDir = vec4( worldVertex - LightPos );
			float lightDistanceSQ = worldLightDir.x*worldLightDir.x + worldLightDir.y*worldLightDir.y + worldLightDir.z*worldLightDir.z;
			float brightness = 1 / lightDistanceSQ;
			worldLightDir = normalize( worldLightDir );

			directColor += ( brightness * directBRDF( LightAmb, LightDiff, LightSpec, worldLightDir, diffuseValue, specularValue, tangentNormal, tangentReflectedEye) );
		}

		//directional lights
		for (int i=0;i<nbDirectionalLights;i++){

			//light components
			vec4 LightPos = directionalLights[(i*4)+LIGHT_POSITION];
			vec4 LightAmb = directionalLights[(i*4)+LIGHT_AMBIENT];
			vec4 LightDiff = directionalLights[(i*4)+LIGHT_DIFFUSE];
			vec4 LightSpec = directionalLights[(i*4)+LIGHT_SPECULAR];

			vec4 worldLightDir = normalize( LightPos );

			directColor += directBRDF( LightAmb, LightDiff, LightSpec, worldLightDir, diffuseValue, specularValue, tangentNormal, tangentReflectedEye);
		}

		//indirect lighting
		vec4 indirectDiffuse = texture(illuminanceMap, worldReflectedEye.xyz);
		vec4 indirectSpecular = vec4(0,0,0,1);
		if (roughnessValue < 0.1){
			indirectSpecular = texture(environmentMap, worldReflectedEye.xyz);
		} else if (roughnessValue < 0.3){
			indirectSpecular = texture(environmentMapLOD1, worldReflectedEye.xyz);
		} else if (roughnessValue < 0.6){
			indirectSpecular = texture(environmentMapLOD2, worldReflectedEye.xyz);
		} else if (roughnessValue < 0.8){
			indirectSpecular = texture(environmentMapLOD3, worldReflectedEye.xyz);
		} else {
			indirectSpecular = texture(illuminanceMap, worldReflectedEye.xyz);
		}

		float freznel = pow((1.0-dot(-tangentEyeDirection, tangentNormal.xyz)), 2);
		//blend indirect light types
		indirectColor += (diffuseValue - vec4(freznel)) * indirectDiffuse;
		indirectColor += (specularValue + vec4(freznel)) * indirectSpecular;

		// finalColor += directColor;
		finalColor += indirectColor;

		finalColor.a = alphaValue + (alphaValue * freznel);

		finalColor = vectorCap(finalColor, 0.99);
	}

	if (mode == MODE_UNLIT){
		finalColor = diffuseValue;
		finalColor = vectorCap(finalColor, 0.99);
	}

	if (mode == MODE_EMIT){
		finalColor = diffuseValue;
	}

	//final output
	outputColor = finalColor;
}
` + "\x00"
