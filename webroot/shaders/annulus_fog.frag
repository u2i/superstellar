varying vec2 vTextureCoord;

uniform vec2 worldCoordinates;
uniform vec2 worldSize;

uniform mat3 magicMatrix;

void main() {
	vec3 mapCoord = vec3(vTextureCoord, 1.0) * magicMatrix;

	float curX = worldCoordinates.x + mapCoord.x * 800.0 - 400.0;
	float curY = -worldCoordinates.y + mapCoord.y * 600.0 - 300.0;

	float dist = distance(vec2(curX, curY), vec2(0.0, 0.0));

	float alpha = smoothstep(worldSize.x, worldSize.y, dist) * 0.6;

	gl_FragColor = vec4(0.603 * alpha, 0.192 * alpha, 0.992 * alpha, alpha);
}
