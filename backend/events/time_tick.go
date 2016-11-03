package events

type TimeTick struct {
	FrameId uint32
}

type TimeTickListener interface {
	HandleTimeTick(*TimeTick)
}
