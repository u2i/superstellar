package backend

import "fmt"

type Move struct {
	ClientID  string	`json:"client_id"`
	Direction string	`json:"direction"`
}

func (self *Move) String() string {
	return fmt.Sprintf("%d moves to %s", self.ClientID, self.Direction)
}
