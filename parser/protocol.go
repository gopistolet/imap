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

type SelectCmd struct {
	Mailbox string
}

type ExamineCmd struct {
	Mailbox string
}

type CreateCmd struct {
	Mailbox string
}

type DeleteCmd struct {
	Mailbox string
}
