#version 330

#vert
in vec2 vert;
in vec2 texCoord;
out vec2 fragTexCoord;
#endvert

#frag
uniform sampler2D diffuseMap;
in vec2 fragTexCoord;
out vec4 outputColor;

uniform float weight[5] = float[] (0.227027, 0.1945946, 0.1216216, 0.054054, 0.016216);
#endfrag

void main() {
  #vert
  gl_Position = vec4(vert, 0.0, 1.0);
  fragTexCoord = texCoord;
  #endvert

  #frag
  vec2 tex_offset = 1.0 / textureSize(diffuseMap, 0);
  vec3 result = texture(diffuseMap, fragTexCoord).rgb * weight[0];
  for (int i = 1; i < 5.; ++i) {
    result += texture(diffuseMap, fragTexCoord + vec2(tex_offset.x * i, 0.0)).rgb * weight[i];
    result += texture(diffuseMap, fragTexCoord - vec2(tex_offset.x * i, 0.0)).rgb * weight[i];
  }
  for (int i = 1; i < 5; ++i) {
    result += texture(diffuseMap, fragTexCoord + vec2(0.0, tex_offset.y * i)).rgb * weight[i];
    result += texture(diffuseMap, fragTexCoord - vec2(0.0, tex_offset.y * i)).rgb * weight[i];
  }
  outputColor = vec4(result, 1.0);
  #endfrag
}