package webgl

import (
	"fmt"
	
	"github.com/gopherjs/webgl"
	"github.com/gopherjs/gopherjs/js"
)

var mainVert = `
uniform vec4 worldCamPos;
uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;
uniform mat4 modelNormal;

attribute vec3 vert;
attribute vec3 vertNormal;
attribute vec2 vertTexCoord;
attribute vec4 color;

varying vec2 fragTexCoord;
varying mat3 TBNMatrix;
varying mat3 inverseTBNMatrix;
varying vec4 worldVertex;
varying vec3 worldNormal;
varying vec4 fragColor;
varying vec3 tangentEyeDirection;

mat3 inverse(mat3 m) {
  float a00 = m[0][0], a01 = m[0][1], a02 = m[0][2];
  float a10 = m[1][0], a11 = m[1][1], a12 = m[1][2];
  float a20 = m[2][0], a21 = m[2][1], a22 = m[2][2];

  float b01 = a22 * a11 - a12 * a21;
  float b11 = -a22 * a10 + a12 * a20;
  float b21 = a21 * a10 - a11 * a20;

  float det = a00 * b01 + a01 * b11 + a02 * b21;

  return mat3(b01, (-a22 * a01 + a02 * a21), (a12 * a01 - a02 * a11),
              b11, (a22 * a00 - a02 * a20), (-a12 * a00 + a02 * a10),
              b21, (-a21 * a00 + a01 * a20), (a11 * a00 - a01 * a10)) / det;
}

void main() {
	worldVertex = model * vec4(vert, 1);
	gl_Position = projection * camera * worldVertex;
	worldNormal = (modelNormal * vec4(vertNormal,1)).xyz;
	vec3 tangent = cross(vertNormal, vertNormal + vec3(-1));
	vec3 bitangent = cross(vertNormal, tangent);
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
`

var mainFrag = `
#define MODE_UNLIT 0
#define MODE_LIT 1
#define MODE_EMIT 2

precision highp float;

uniform int mode;
uniform bool useVertexColor;

uniform vec4 directionalLightPos;
uniform vec4 directionalLightAmb;
uniform vec4 directionalLightDiff;
uniform vec4 directionalLightSpec;

//material
uniform sampler2D diffuse;
uniform sampler2D normal;
uniform sampler2D specular;
uniform sampler2D roughness;
uniform samplerCube environmentMap;
uniform samplerCube environmentMapLOD1;
uniform samplerCube environmentMapLOD2;
uniform samplerCube environmentMapLOD3;
uniform samplerCube illuminanceMap;

uniform vec4 worldCamPos;
uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;

varying vec2 fragTexCoord;
varying mat3 TBNMatrix;
varying mat3 inverseTBNMatrix;
varying vec4 worldVertex;
varying vec3 worldNormal;
varying vec4 fragColor;
varying vec3 tangentEyeDirection;

vec4 vectorCap( vec4 vector, float cap ){
	if (vector.r > cap) {vector.r = cap;}
	if (vector.g > cap) {vector.g = cap;}
	if (vector.b > cap) {vector.b = cap;}
	return vector;
}

vec4 directBRDF( vec4 LightAmb, vec4 LightDiff, vec4 LightSpec, vec4 LightDir, vec4 albedoValue, vec4 specularValue, vec4 tangentNormal, vec4 tangentReflectedEye){
	vec3 tangentLightDirection = LightDir.xyz * TBNMatrix;
	tangentLightDirection = normalize( tangentLightDirection );

	vec4 ambientOut = albedoValue * LightAmb;

	float diffuseMultiplier = max(0.0, dot(tangentNormal.xyz, -tangentLightDirection));
	vec4 diffuseOut = albedoValue * diffuseMultiplier * LightDiff;

	float specularMultiplier = pow( max(0.0, dot( tangentReflectedEye.xyz, -tangentLightDirection)), 2.0);
	vec4 specularOut = vec4( specularValue.rgb, 1 ) * specularMultiplier * LightSpec;

	return ambientOut + diffuseOut + specularOut;
}

void main() {
	float textureX = fragTexCoord.x;
	float textureY = fragTexCoord.y;
	if (fragTexCoord.x < 0.0) {textureX = textureX + 1.0;}
	if (fragTexCoord.y < 0.0) {textureY = textureY + 1.0;}
	vec2 textCoord = vec2(textureX, textureY);

	vec4 albedoValue = texture2D(diffuse, textCoord) * fragColor;
	if (useVertexColor) {
		albedoValue = fragColor;
	}
	vec4 normalValue = texture2D(normal, textCoord);
	vec4 specularValue = texture2D(specular, textCoord);
	vec4 roughnessValue = texture2D(roughness, textCoord);
	float roughnessMagnitude = roughnessValue.r;
	float metalness = roughnessValue.g;
	float alphaValue = albedoValue.a;

	vec4 finalColor = vec4(0,0,0,1);
	vec4 directColor = vec4(0,0,0,1);
	vec4 indirectColor = vec4(0,0,0,1);

	if( mode == MODE_LIT ){

		//Normal calculations
		vec4 tangentNormal = normalValue * 2.0 - 1.0;
		if( abs(tangentNormal.x) < 0.1 && abs(tangentNormal.y) < 0.1 && abs(tangentNormal.z) < 0.1 ){
			tangentNormal = vec4(0,0,1,1);
		}

		//reflected eye
		vec4 tangentReflectedEye = vec4( reflect( tangentEyeDirection, tangentNormal.xyz ), 1);
		vec4 worldReflectedEye = vec4( tangentReflectedEye.xyz * inverseTBNMatrix , 1);

		//directional light

		//light components
		vec4 LightPos = directionalLightPos;
		vec4 LightAmb = directionalLightAmb;
		vec4 LightDiff = directionalLightDiff;
		vec4 LightSpec = directionalLightSpec;

		vec4 worldLightDir = normalize( directionalLightPos );

		directColor += directBRDF( directionalLightAmb, directionalLightDiff, directionalLightSpec, worldLightDir, albedoValue, specularValue, tangentNormal, tangentReflectedEye);

		//indirect lighting
		vec4 indirectDiffuse = textureCube(illuminanceMap, worldReflectedEye.xyz);
		vec4 indirectSpecular = vec4(0,0,0,1);
		if (roughnessMagnitude < 0.1){
			indirectSpecular = textureCube(environmentMap, worldReflectedEye.xyz);
		} else if (roughnessMagnitude < 0.3){
			indirectSpecular = textureCube(environmentMapLOD1, worldReflectedEye.xyz);
		} else if (roughnessMagnitude < 0.7){
			indirectSpecular = textureCube(environmentMapLOD2, worldReflectedEye.xyz);
		} else {
			indirectSpecular = textureCube(environmentMapLOD3, worldReflectedEye.xyz);
		}

		//freznel effect
		float f = (1.0-dot(-tangentEyeDirection, tangentNormal.xyz));
		float reflectivity = max(metalness, f*f);
		//blend indirect light types
		indirectColor += (1.0-metalness) * albedoValue  * indirectDiffuse;
		indirectColor += reflectivity * specularValue * indirectSpecular;

		finalColor += (1.0-metalness) * directColor;
		finalColor += indirectColor;

		finalColor.a = alphaValue + (alphaValue * reflectivity);

		finalColor = vectorCap(finalColor, 0.99);
	}

	if (mode == MODE_UNLIT){
		finalColor = albedoValue;
		finalColor = vectorCap(finalColor, 0.99);
	}

	if (mode == MODE_EMIT){
		finalColor = albedoValue;
	}

	//final output
	gl_FragColor = finalColor;
}
`

func compileShader(gl *webgl.Context, source string, shaderType int) *js.Object {
	shader := gl.CreateShader(shaderType);
	gl.ShaderSource(shader, source);
	gl.CompileShader(shader);
	
  status := gl.GetShaderParameter(shader, gl.COMPILE_STATUS);
	if !status.Bool() {
		panic(fmt.Sprintf("failed to compile shader: %v", gl.GetShaderInfoLog(shader)))
  }

	return shader
}