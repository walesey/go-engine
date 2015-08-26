#version 330

#define MAX_LIGHTS 8

#define LIGHT_POSITION 0
#define LIGHT_AMBIENT 1
#define LIGHT_DIFFUSE 2
#define LIGHT_SPECULAR 3

uniform vec4 lights[ MAX_LIGHTS * 4 ];

//material
uniform sampler2D diffuse;
uniform sampler2D normal;
uniform sampler2D specular;
uniform sampler2D roughness;
uniform samplerCube environmentMap;

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;

in mat3 TBNMatrix;
in mat3 inverseTBNMatrix;
in vec4 worldCamPos;
in vec4 worldVertex;
in vec3 worldNormal;
in vec2 fragTexCoord;

out vec4 outputColor;

void main() {
	vec4 diffuseValue = texture(diffuse, fragTexCoord);
	vec4 normalValue = texture(normal, fragTexCoord);
	vec4 specularValue = texture(specular, fragTexCoord);
	vec4 roughnessValue = texture(roughness, fragTexCoord);
	vec4 finalColor = vec4(0,0,0,1);

	//Normal calculations
 	vec3 normalMapValue = vec3( normalValue.rgb ) * 2 - 1;
 	if( abs(normalMapValue.x) < 0.1 && abs(normalMapValue.y) < 0.1 && abs(normalMapValue.z) < 0.1 ){
 		normalMapValue = vec3(0,0,1);
 	}

	//eye 
	vec4 worldEyeDirection = vec4( worldVertex - worldCamPos );
	vec3 tangentEyeDirection = normalize( worldEyeDirection.xyz * TBNMatrix );
	vec4 tangentReflectedEye = vec4( reflect( tangentEyeDirection, normalMapValue ), 1);

	//lights
	for (int i=0;i<MAX_LIGHTS;i++){

		//light components
		vec4 LightPos = lights[(i*4)+LIGHT_POSITION];
		vec4 LightAmb = lights[(i*4)+LIGHT_AMBIENT];
		vec4 LightDiff = lights[(i*4)+LIGHT_DIFFUSE];
		vec4 LightSpec = lights[(i*4)+LIGHT_SPECULAR];

		//point light source
		vec4 worldLightDir = vec4( LightPos - worldVertex );
		float lightDistanceSQ = worldLightDir.x*worldLightDir.x + worldLightDir.y*worldLightDir.y + worldLightDir.z*worldLightDir.z;
		float illuminanceMultiplier = 1 / lightDistanceSQ;
		
		//directional light source (only index 0)
		if( i == 0 ){
			worldLightDir = normalize( LightPos );
			illuminanceMultiplier = 1;
		}

		//tangent space
		vec3 tangentLightDirection = normalize( worldLightDir.xyz * TBNMatrix );

		//ambient component
		vec4 ambientOut = diffuseValue * LightAmb;
		//diffuse component
	 	float diffuseMultiplier = max(0.0, dot(normalMapValue, tangentLightDirection));
		vec4 diffuseOut = diffuseValue * diffuseMultiplier * LightDiff;
		//specular component
	 	float specularMultiplier = max(0.0, pow( dot(tangentReflectedEye.xyz, tangentLightDirection), 2.0));
		vec4 specularOut = vec4( specularValue.rgb, 1 ) * specularMultiplier * LightSpec;

		finalColor += (( ambientOut + diffuseOut + specularOut ) * illuminanceMultiplier);
   	}

   	//indirect lighting
   	vec4 worldReflectedEye = vec4( tangentReflectedEye.xyz * inverseTBNMatrix , 1);
   	worldReflectedEye = worldReflectedEye + vec4( ( noise3(1) * roughnessValue.xyz ), 1 );
	finalColor = texture(environmentMap, worldReflectedEye.xyz) * specularValue.a;

	//final output
	outputColor = finalColor;
}