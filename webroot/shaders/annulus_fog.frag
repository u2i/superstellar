

// Star Nest by Pablo RomÃ¡n Andrioli

// This content is under the MIT License.

#define iterations 17
#define formuparam 0.53

#define volsteps 20
#define stepsize 0.1

#define zoom   0.800
#define tile   0.850
#define speed  0.010

#define brightness 0.0015
#define darkmatter 0.300
#define distfading 0.730
#define saturation 0.850

varying vec2 vTextureCoord;

uniform vec2 worldCoordinates;
uniform vec2 worldSize;
uniform vec4 dimensions;

uniform mat3 magicMatrix;

void main() {
	vec3 mapCoord = vec3(vTextureCoord, 1.0) * magicMatrix;

	float curX = worldCoordinates.x + mapCoord.x * dimensions.x - dimensions.x / 2.0;
	float curY = -worldCoordinates.y + mapCoord.y * dimensions.y - dimensions.y / 2.0;

	float dist = distance(vec2(curX, curY), vec2(0.0, 0.0));

	float alpha = smoothstep(worldSize.x, worldSize.y, dist) * 0.6;

	gl_FragColor = vec4(0.303 * alpha, 0.192 * alpha, 0.992 * alpha, alpha);
}
