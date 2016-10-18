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
			command = LoginCmd{
				Username: string(lexCommand.Arguments[0]),
				Password: string(lexCommand.Arguments[1]),
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
	default:
		{
			err = errors.New("Unknown command")
		}
	}

	return

}
