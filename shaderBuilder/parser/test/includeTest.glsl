#include "includeTest2.glsl"

#frag
uniform #lookup 10 lookupTable ((i + 1) / 100.0);
uniform #lookup 10 lookupTable2 ((i + 1) / -100.0);

vec4 fragFunc() {
  output = vec4(0.5);
  for (int i = 0; i < 10; ++i) {
    output += vec4(lookupTable[i]);
  }
  return vec4(0.5);
}
#endfrag

#vert
vec4 vertFunc() {
  return vec4(1.0);
}
#endvert