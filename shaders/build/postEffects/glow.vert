#version 330

in vec2 vert;
in vec2 texCoord;
out vec2 fragTexCoord;

void main() {
  
  gl_Position = vec4(vert, 0.0, 1.0);
  fragTexCoord = texCoord;

}
