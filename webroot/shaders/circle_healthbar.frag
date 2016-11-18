#define PI 3.1415926535897932384626433832795
precision mediump float;

varying vec2 vTextureCoord;

uniform vec2 hps;
uniform mat3 magicMatrix;

// colors
const vec3 healthColor = vec3(0.5, 1.0, 0.5);
const float healthAlpha = 0.4;
const float lostHealthAlpha = 0.1;

// circle constants
const vec2 center = vec2(0.5, 0.5);
const float radius = 0.45;
const float epsilon = 0.030;
const float angleFraction = PI;
// hp bar consts
const float hpPerBucket = 100.0; // because of circle mirroring each bucket shows half of that value

void main() {
    float hp = hps.x;
    float maxHp = hps.y;
    // hp bar values calculated from HP
    float bucketsNo = ceil(maxHp / hpPerBucket);
    float bucketWidth = 1.0 / bucketsNo;
    float separatorWidth = 0.08 * bucketWidth;
    float separatorHalfWidth = separatorWidth / 2.0;
    float bucketInteriorWidth = bucketWidth - separatorWidth;
    float health = hp / maxHp;

    vec3 mapCoord = vec3(vTextureCoord, 1.0) * magicMatrix;
    vec2 u = mapCoord.xy;
    vec2 translated_u = vec2(abs(u.x - center.x), u.y - center.y);

    float circleFilter = step(abs(length(translated_u) - radius), epsilon);

    // translate x, y to spherical coords
    float vectorAngle = atan(translated_u.y / translated_u.x);
    float omega = PI / 2.0 - vectorAngle;
    float angleFilled = 1.0 / angleFraction * omega;

    // kudos to Mateja
    float x = angleFilled;
    float bucketsBehind = floor((x + separatorHalfWidth) / bucketWidth);
    float logicalX = ((x - separatorHalfWidth) - (bucketsBehind * separatorWidth)) / (bucketsNo * bucketInteriorWidth);
    float logicalXHealth = logicalX * (bucketsNo * hpPerBucket / maxHp);
    float hpBarFunction = step(logicalXHealth, health);
    float alphaFunction = step(separatorHalfWidth, abs((bucketsBehind * bucketWidth) - x));

    float lostHealthOverlay = alphaFunction * lostHealthAlpha;

    float alpha = circleFilter * healthAlpha * hpBarFunction * alphaFunction;
    vec3 color = (circleFilter * alpha + circleFilter * lostHealthOverlay) * healthColor;

    gl_FragColor = vec4(color, alpha);
}
