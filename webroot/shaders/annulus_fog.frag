varying vec2 vTextureCoord;

uniform vec2 worldCoordinates;
uniform vec2 worldBoundarySize;
uniform vec2 windowSize;

uniform mat3 magicMatrix;

void main() {
	vec3 mapCoord = vec3(vTextureCoord, 1.0) * magicMatrix;

	float curX = worldCoordinates.x + mapCoord.x * windowSize.x - windowSize.x / 2.0;
	float curY = -worldCoordinates.y + mapCoord.y * windowSize.y - windowSize.y / 2.0;

	float dist = distance(vec2(curX, curY), vec2(0.0, 0.0));

	float alpha = smoothstep(worldBoundarySize.x, worldBoundarySize.y, dist) * 0.6;

	gl_FragColor = vec4(0.303 * alpha, 0.192 * alpha, 0.992 * alpha, alpha);
}
