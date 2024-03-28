package error

type Code int

const (
	Unknown Code = 0
	Success Code = 1000 + iota
	InvalidRequest
	Internal

	PartyNotFound
	CouldNotCreateParty
)

func (c Code) String() string {
	if c == 0 {
		return "Unknown"
	}

	return [...]string{
		"Success",
		"Invalid Request",
		"Internal Error",
		"Party Not Found",
		"Could Not Create Party",
	}[c-1000]
}
