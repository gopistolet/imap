package parser

import (
	"errors"
)

// parseLine parses a single line and returns the matching IMAP command
func parseLine(line string) (command Cmd, err error) {

	lexCommand, err := lexLine(line)
	if err != nil {
		return
	}

	switch lexCommand.Name {
	case "STARTTLS":
		{
			if len(lexCommand.Arguments) != 0 {
				err = errors.New("Parser: expected no arguments for STARTTLS command")
				return
			}
			command = StarttlsCmd{}
		}
	case "LOGIN":
		{
			/*
				login    = "LOGIN" SP userid SP password
				userid   = astring
				password = astring
			*/
			if len(lexCommand.Arguments) != 2 {
				err = errors.New("Parser: expected two arguments (username, password) for LOGIN command")
				return
			}
			if !isAString(lexCommand.Arguments[0]) {
				err = errors.New("Parser: expected first argument (userid) to be astring")
				return
			}
			if !isAString(lexCommand.Arguments[1]) {
				err = errors.New("Parser: expected second argument (password) to be astring")
				return
			}
			command = LoginCmd{
				Username: lexCommand.Arguments[0],
				Password: lexCommand.Arguments[1],
			}
		}
	case "LOGOUT":
		{
			if len(lexCommand.Arguments) != 0 {
				err = errors.New("Parser: expected no arguments for LOGOUT command")
				return
			}
			command = LogoutCmd{}
		}
	case "CAPABILITY":
		{
			if len(lexCommand.Arguments) != 0 {
				err = errors.New("Parser: expected no arguments for CAPABILITY command")
				return
			}
			command = CapabilityCmd{}
		}
	case "NOOP":
		{
			if len(lexCommand.Arguments) != 0 {
				err = errors.New("Parser: expected no arguments for NOOP command")
				return
			}
			command = NoopCmd{}
		}
	case "AUTHENTICATE":
		{
			/*
				authenticate    = "AUTHENTICATE" SP auth-type *(CRLF base64)
				auth-type       = atom
									; Defined by [SASL]
			*/
			if len(lexCommand.Arguments) != 1 {
				err = errors.New("Parser: expected 1 argument (authentication mechanism name) for AUTHENTICATE command")
				return
			}
			if !isAtom(lexCommand.Arguments[0]) {
				err = errors.New("Parser: expected first argument (authentication mechanism name) to be atom")
				return
			}

			command = AuthenticateCmd{
				Mechanism: lexCommand.Arguments[0],
			}
		}
	default:
		{
			err = errors.New("Unknown command")
		}
	}

	return

}
