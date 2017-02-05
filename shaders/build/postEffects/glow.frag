#version 330

uniform sampler2D tex0;
uniform sampler2D tex1;
in vec2 fragTexCoord;
out vec4 outputColor;

uniform float weight[5] = float[] (0.227027, 0.1945946, 0.1216216, 0.054054, 0.016216);

void main() {

  vec2 tex_offset = 1.0 / textureSize(tex0, 0);
  vec3 result = texture(tex1, fragTexCoord).rgb;
  result += texture(tex0, fragTexCoord).rgb * weight[0];
  for (int i = 1; i < 5; ++i) {
    result += texture(tex0, fragTexCoord + vec2(tex_offset.x * i, 0.0)).rgb * weight[i];
    result += texture(tex0, fragTexCoord - vec2(tex_offset.x * i, 0.0)).rgb * weight[i];
  }
  for (int i = 1; i < 5; ++i) {
    result += texture(tex0, fragTexCoord + vec2(0.0, tex_offset.y * i)).rgb * weight[i];
    result += texture(tex0, fragTexCoord - vec2(0.0, tex_offset.y * i)).rgb * weight[i];
  }
  outputColor = vec4(result, 1.0);
  
}
