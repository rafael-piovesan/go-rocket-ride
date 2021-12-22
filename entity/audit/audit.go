package audit

type Action string

const (
	ActionCreateRide Action = "CREATE_RIDE"
)

func (a Action) String() string {
	return string(a)
}

type ResourceType string

const (
	ResourceTypeRide ResourceType = "RIDE"
)

func (a ResourceType) String() string {
	return string(a)
}
