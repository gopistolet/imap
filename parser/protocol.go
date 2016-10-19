package parser

type Cmd interface {
}

type LoginCmd struct {
	Username string
	Password string
}

type LogoutCmd struct {
}

type CapabilityCmd struct {
}

type NoopCmd struct {
}

type StarttlsCmd struct {
}

type AuthenticateCmd struct {
	Mechanism string
}
