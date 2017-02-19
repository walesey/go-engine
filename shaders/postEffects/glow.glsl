#version 330

#vert
in vec2 vert;
in vec2 texCoord;
out vec2 fragTexCoord;
#endvert

#frag
uniform sampler2D tex0;
uniform sampler2D tex1;
in vec2 fragTexCoord;
out vec4 outputColor;

uniform float sample = 2.0;
uniform float iterations = 8;
uniform #lookup 100 gaussian (2.718 ^ ( -(i*0.02 + 1.2) ^ 2.0 ));
#endfrag

void main() {
  #vert
  gl_Position = vec4(vert, 0.0, 1.0);
  fragTexCoord = texCoord;
  #endvert

  #frag
  vec2 tex_offset = sample / textureSize(tex1, 0);
  vec3 result = texture(tex0, fragTexCoord).rgb;
  result += texture(tex1, fragTexCoord).rgb * gaussian[0];
  for (int i = 1; i < iterations; ++i) {
    result += texture(tex1, fragTexCoord + vec2(tex_offset.x * i, 0.0)).rgb * gaussian[int((100*i)/iterations)]*sample;
    result += texture(tex1, fragTexCoord - vec2(tex_offset.x * i, 0.0)).rgb * gaussian[int((100*i)/iterations)]*sample;
  }
  for (int i = 1; i < iterations; ++i) {
    result += texture(tex1, fragTexCoord + vec2(0.0, tex_offset.y * i)).rgb * gaussian[int((100*i)/iterations)]*sample;
    result += texture(tex1, fragTexCoord - vec2(0.0, tex_offset.y * i)).rgb * gaussian[int((100*i)/iterations)]*sample;
  }
  outputColor = vec4(result, 1.0);
  #endfrag
}