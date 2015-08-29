#version 330

#define MAX_LIGHTS 8

#define LIGHT_POSITION 0
#define LIGHT_AMBIENT 1
#define LIGHT_DIFFUSE 2
#define LIGHT_SPECULAR 3

#define MODE_UNLIT 0
#define MODE_LIT 1

uniform int mode;

uniform vec4 lights[ MAX_LIGHTS * 4 ];
uniform vec4 directionalLights[ MAX_LIGHTS * 4 ];

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

vec4 directLight( vec4 LightAmb, vec4 LightDiff, vec4 LightSpec, vec4 LightDir, 
				vec4 diffuseValue, vec4 specularValue, vec4 normalMapValue, 
				vec4 tangentReflectedEye, float illuminanceMultiplier){
	//tangent space
	vec3 tangentLightDirection = LightDir.xyz * TBNMatrix;

	//ambient component
	vec4 ambientOut = diffuseValue * LightAmb;
	//diffuse component
 	float diffuseMultiplier = max(0.0, dot(normalMapValue.xyz, -tangentLightDirection));
	vec4 diffuseOut = diffuseValue * diffuseMultiplier * LightDiff;
	//specular component
 	float specularMultiplier = max(0.0, pow( dot(tangentReflectedEye.xyz, -tangentLightDirection), 3.0));
	vec4 specularOut = vec4( specularValue.rgb, 1 ) * specularMultiplier * LightSpec;

	return (( ambientOut + diffuseOut + specularOut ) * illuminanceMultiplier);
}

void main() {
	vec4 diffuseValue = texture(diffuse, fragTexCoord);
	vec4 normalValue = texture(normal, fragTexCoord);
	vec4 specularValue = texture(specular, fragTexCoord);
	vec4 roughnessValue = texture(roughness, fragTexCoord);
	vec4 finalColor = vec4(0,0,0,1);
	vec4 indirectValue = vec4(0,0,0,1);


 	if( mode == MODE_LIT ){

		//Normal calculations
	 	vec4 normalMapValue = normalValue * 2 - 1;
	 	if( abs(normalMapValue.x) < 0.1 && abs(normalMapValue.y) < 0.1 && abs(normalMapValue.z) < 0.1 ){
	 		normalMapValue = vec4(0,0,1,1);
	 	}

		//eye 
		vec4 worldEyeDirection = vec4( worldVertex - worldCamPos );
		vec3 tangentEyeDirection = normalize( worldEyeDirection.xyz * TBNMatrix );
		vec4 tangentReflectedEye = vec4( reflect( tangentEyeDirection, normalMapValue.xyz ), 1);
	   	vec4 worldReflectedEye = vec4( tangentReflectedEye.xyz * inverseTBNMatrix , 1);

		//point lights
		for (int i=0;i<MAX_LIGHTS;i++){

			//light components
			vec4 LightPos = lights[(i*4)+LIGHT_POSITION];
			vec4 LightAmb = lights[(i*4)+LIGHT_AMBIENT];
			vec4 LightDiff = lights[(i*4)+LIGHT_DIFFUSE];
			vec4 LightSpec = lights[(i*4)+LIGHT_SPECULAR];

			//point light source
			vec4 worldLightDir = vec4( worldVertex - LightPos );
			float lightDistanceSQ = worldLightDir.x*worldLightDir.x + worldLightDir.y*worldLightDir.y + worldLightDir.z*worldLightDir.z;
			float illuminanceMultiplier = 1 / lightDistanceSQ;
			worldLightDir = normalize( worldLightDir );

			finalColor += directLight( 
				LightAmb, LightDiff, LightSpec, worldLightDir,
				diffuseValue, specularValue, normalMapValue, tangentReflectedEye, illuminanceMultiplier );
	   	}

		//directional light
	   	for (int i=0;i<MAX_LIGHTS;i++){

			//light components
			vec4 LightPos = directionalLights[(i*4)+LIGHT_POSITION];
			vec4 LightAmb = directionalLights[(i*4)+LIGHT_AMBIENT];
			vec4 LightDiff = directionalLights[(i*4)+LIGHT_DIFFUSE];
			vec4 LightSpec = directionalLights[(i*4)+LIGHT_SPECULAR];
			
			vec4 worldLightDir = normalize( LightPos );

			finalColor += directLight( 
				LightAmb, LightDiff, LightSpec, worldLightDir,
				diffuseValue, specularValue, normalMapValue, tangentReflectedEye, 1 );
		}

	   	//indirect lighting
	   	//TODO: roughness - gaussian blur the cubemap
		indirectValue += texture(environmentMap, -worldReflectedEye.xyz);

		//freznel effect
		float ratio = 1.0-specularValue.a;
		ratio = min( 1.0, ratio + 0.7*pow(1-dot(-tangentEyeDirection, normalMapValue.xyz), 3) );

		//blend direct/indirect
		finalColor = vec4(min(1.0,finalColor.r), min(1.0,finalColor.g), min(1.0,finalColor.b), 1);
		finalColor = (1.0-ratio)*finalColor + ratio*diffuseValue*indirectValue;
	}
	
	if( mode == MODE_UNLIT ){
		finalColor = diffuseValue;
	}

	//final output
	outputColor = finalColor;
}