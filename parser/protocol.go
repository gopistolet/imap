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

type AuthenticatedStateCmd interface {
	GetMailbox() string
}

type SelectCmd struct {
	Mailbox string
}

func (cmd SelectCmd) GetMailbox() string {
	return cmd.Mailbox
}

type ExamineCmd struct {
	Mailbox string
}

func (cmd ExamineCmd) GetMailbox() string {
	return cmd.Mailbox
}

type CreateCmd struct {
	Mailbox string
}

func (cmd CreateCmd) GetMailbox() string {
	return cmd.Mailbox
}

type DeleteCmd struct {
	Mailbox string
}

func (cmd DeleteCmd) GetMailbox() string {
	return cmd.Mailbox
}
