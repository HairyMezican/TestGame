package player

type Group interface {
	Playable(user string) (bool, string)
	Count() int
	JoinMethods(user string) map[string]JoinMethod
}

type JoinMethod interface {
	Name() string
}
