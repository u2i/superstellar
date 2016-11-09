package events

type UserDied struct {
	ClientID	uint32
	KilledBy	uint32
}
