precision mediump float;
#define PI 3.1415926535897932384626433832795

varying vec2 vTextureCoord;

uniform vec2 hps;
uniform mat3 magicMatrix;
uniform vec3 basicColor;

// colors
const float healthAlpha = 0.4;
const float lostHealthAlpha = 0.2;

// circle constants
const vec2 center = vec2(0.5, 0.5);
const float radius = 0.45;
const float epsilon = 0.035;
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
    float health = hp / (bucketsNo * hpPerBucket);

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
    float hpBarFunction = step(logicalX, health);
    float alphaFunction = healthAlpha * step(separatorHalfWidth, abs((bucketsBehind * bucketWidth) - x)) * step(logicalX, maxHp/(bucketsNo * hpPerBucket));
    float lostHealthOverlay = alphaFunction * lostHealthAlpha;

    float colorFunction = hpBarFunction * alphaFunction + lostHealthOverlay;
    gl_FragColor = vec4(circleFilter * basicColor * colorFunction, circleFilter * hpBarFunction * alphaFunction);
}
