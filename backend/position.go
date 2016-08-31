package backend

import "fmt"

type Position struct {
	ClientID string  `json:"client_id"`
	X        int     `json:"x"`
	Y        int     `json:"y"`
	Angle    float32 `json:"angle"`
}

func (self *Position) String() string {
	return fmt.Sprintf("%s is at (%d, %d)", self.ClientID, self.X, self.Y)
}
