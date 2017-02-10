#version 330

float import1() {
	return 1.0;
}

vec4 vertFunc() {
  return vec4(1.0);
}

void main() {
  gl_Position = vertFunc();
}

