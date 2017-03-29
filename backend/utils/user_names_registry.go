package utils

type UserNamesRegistry struct {
	userNames map[uint32]string
}

func NewUserNameRegistry() *UserNamesRegistry {
	return &UserNamesRegistry{
		userNames: make(map[uint32]string),
	}
}

func (registry *UserNamesRegistry) AddUserName(id uint32, userName string) {
	registry.userNames[id] = userName
}

func (registry *UserNamesRegistry) GetUserName(id uint32) string {
	return registry.userNames[id]
}
