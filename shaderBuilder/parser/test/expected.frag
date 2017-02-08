#version 330

float import1() {
	return 1.0;
}

uniform float lookupTable[10] = float[] (0.01, 0.02, 0.03, 0.04, 0.05, 0.06, 0.07, 0.08, 0.09, 0.1);

vec4 fragFunc() {
  output = vec4(0.5);
  for (int i = 0; i < 10; ++i) {
    output += vec4(lookupTable[i]);
  }
  return vec4(0.5);
}

void main() {
  gl_FragColor = fragFunc();
}
