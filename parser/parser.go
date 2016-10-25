package parser

import (
	"errors"
	"strings"
	"time"
)

// parseLine parses a single line and returns the matching IMAP command
func parseLine(line string) (command Cmd, err error) {

	lexCommand, err := lexLine(line)
	if err != nil {
		return
	}

	switch lexCommand.Name {

	// Client Commands - Any State
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

	// Client Commands - Not Authenticated State
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

	// Client Commands - Authenticated State
	case "SELECT":
		{
			/*
				select  = "SELECT" SP mailbox
			*/
			if len(lexCommand.Arguments) != 1 {
				err = errors.New("Parser: expected 1 argument for SELECT command")
				return
			}
			if !isMailbox(lexCommand.Arguments[0]) {
				err = errors.New("Parser: expected first argument (mailbox) for SELECT to be 'INBOX' or astring")
				return
			}

			command = SelectCmd{
				Mailbox: parseMailbox(lexCommand.Arguments[0]),
			}
		}
	case "EXAMINE":
		{
			/*
				examine = "EXAMINE" SP mailbox
			*/
			if len(lexCommand.Arguments) != 1 {
				err = errors.New("Parser: expected 1 argument for EXAMINE command")
				return
			}
			if !isMailbox(lexCommand.Arguments[0]) {
				err = errors.New("Parser: expected first argument (mailbox) for EXAMINE to be 'INBOX' or astring")
				return
			}

			command = ExamineCmd{
				Mailbox: parseMailbox(lexCommand.Arguments[0]),
			}
		}
	case "CREATE":
		{
			/*
				create = "CREATE" SP mailbox
				          ; Use of INBOX gives a NO error
			*/
			if len(lexCommand.Arguments) != 1 {
				err = errors.New("Parser: expected 1 argument for CREATE command")
				return
			}
			if !isMailbox(lexCommand.Arguments[0]) {
				err = errors.New("Parser: expected first argument (mailbox) for CREATE to be 'INBOX' or astring")
				return
			}

			command = CreateCmd{
				Mailbox: parseMailbox(lexCommand.Arguments[0]),
			}
		}
	case "DELETE":
		{
			/*
				delete = "DELETE" SP mailbox
				          ; Use of INBOX gives a NO error
			*/
			if len(lexCommand.Arguments) != 1 {
				err = errors.New("Parser: expected 1 argument for DELETE command")
				return
			}
			if !isMailbox(lexCommand.Arguments[0]) {
				err = errors.New("Parser: expected first argument (mailbox) for DELETE to be 'INBOX' or astring")
				return
			}

			command = DeleteCmd{
				Mailbox: parseMailbox(lexCommand.Arguments[0]),
			}
		}
	case "RENAME":
		{
			/*
			   rename = "RENAME" SP mailbox SP mailbox
			             ; Use of INBOX as a destination gives a NO error
			*/
			if len(lexCommand.Arguments) != 2 {
				err = errors.New("Parser: expected 2 argument for RENAME command")
				return
			}
			if !isMailbox(lexCommand.Arguments[0]) {
				err = errors.New("Parser: expected first argument (mailbox) for RENAME to be 'INBOX' or astring")
				return
			}
			if !isMailbox(lexCommand.Arguments[1]) {
				err = errors.New("Parser: expected second argument (mailbox) for RENAME to be 'INBOX' or astring")
				return
			}

			command = RenameCmd{
				SourceMailbox:      parseMailbox(lexCommand.Arguments[0]),
				DestinationMailbox: parseMailbox(lexCommand.Arguments[1]),
			}
		}
	case "SUBSCRIBE":
		{
			/*
				subscribe = "SUBSCRIBE" SP mailbox
			*/
			if len(lexCommand.Arguments) != 1 {
				err = errors.New("Parser: expected 1 argument for SUBSCRIBE command")
				return
			}
			if !isMailbox(lexCommand.Arguments[0]) {
				err = errors.New("Parser: expected first argument (mailbox) for SUBSCRIBE to be 'INBOX' or astring")
				return
			}

			command = SubscribeCmd{
				Mailbox: parseMailbox(lexCommand.Arguments[0]),
			}
		}
	case "UNSUBSCRIBE":
		{
			/*
				unsubscribe = "UNSUBSCRIBE" SP mailbox
			*/
			if len(lexCommand.Arguments) != 1 {
				err = errors.New("Parser: expected 1 argument for UNSUBSCRIBE command")
				return
			}
			if !isMailbox(lexCommand.Arguments[0]) {
				err = errors.New("Parser: expected first argument (mailbox) for UNSUBSCRIBE to be 'INBOX' or astring")
				return
			}

			command = UnsubscribeCmd{
				Mailbox: parseMailbox(lexCommand.Arguments[0]),
			}
		}
	case "LIST":
		{
			/*
			   list = "LIST" SP mailbox SP list-mailbox
			*/
			if len(lexCommand.Arguments) != 2 {
				err = errors.New("Parser: expected 2 arguments for LIST command")
				return
			}
			if !isMailbox(lexCommand.Arguments[0]) {
				err = errors.New("Parser: expected first argument (reference) for LIST to be 'INBOX' or astring")
				return
			}
			if !isListMailbox(lexCommand.Arguments[1]) {
				err = errors.New("Parser: expected second argument (mailbox) for LIST to be list-mailbox")
				return
			}

			command = ListCmd{
				Reference: parseMailbox(lexCommand.Arguments[0]),
				Mailbox:   lexCommand.Arguments[1],
			}
		}
	case "LSUB":
		{
			/*
			   lsub = "LSUB" SP mailbox SP list-mailbox
			*/
			if len(lexCommand.Arguments) != 2 {
				err = errors.New("Parser: expected 2 arguments for LSUB command")
				return
			}
			if !isMailbox(lexCommand.Arguments[0]) {
				err = errors.New("Parser: expected first argument (reference) for LSUB to be 'INBOX' or astring")
				return
			}
			if !isListMailbox(lexCommand.Arguments[1]) {
				err = errors.New("Parser: expected second argument (mailbox) for LSUB to be list-mailbox")
				return
			}

			command = LsubCmd{
				Reference: parseMailbox(lexCommand.Arguments[0]),
				Mailbox:   lexCommand.Arguments[1],
			}
		}
	case "STATUS":
		{
			/*
			   status     = "STATUS" SP mailbox SP "(" status-att *(SP status-att) ")"
			   status-att = "MESSAGES" / "RECENT" / "UIDNEXT" / "UIDVALIDITY" / "UNSEEN"
			*/
			if len(lexCommand.Arguments) < 1 {
				err = errors.New("Parser: expected mailbox argument for STATUS command")
				return
			}
			if !isMailbox(lexCommand.Arguments[0]) {
				err = errors.New("Parser: expected first argument (mailbox) for STATUS to be 'INBOX' or astring")
				return
			}
			if len(lexCommand.Arguments) < 2 {
				err = errors.New("Parser: expected status-att list for STATUS command")
				return
			}
			firstAttr := &lexCommand.Arguments[1]
			lastAttr := &lexCommand.Arguments[len(lexCommand.Arguments)-1]
			if (*firstAttr)[0] != '(' || (*lastAttr)[len(*lastAttr)-1] != ')' {
				err = errors.New("Parser: expected status-att list for STATUS command. Didn't find delimiter")
				return
			}
			*firstAttr = strings.TrimPrefix(*firstAttr, "(")
			*lastAttr = strings.TrimSuffix(*lastAttr, ")")

			for _, attr := range lexCommand.Arguments[1:] {
				switch attr {
				case "MESSAGES", "RECENT", "UIDNEXT", "UIDVALIDITY", "UNSEEN":
					{
						continue
					}
				default:
					{
						err = errors.New("Parser: unknown status-att for STATUS command: " + attr)
						return
					}

				}
			}

			command = StatusCmd{
				Mailbox:          parseMailbox(lexCommand.Arguments[0]),
				StatusAttributes: lexCommand.Arguments[1:],
			}

		}
	case "APPEND":
		{
			/*
				append          = "APPEND" SP mailbox [SP flag-list] [SP date-time] SP literal
				flag-list       = "(" [flag *(SP flag)] ")"
				flag            = "\Answered" / "\Flagged" / "\Deleted" /
				                  "\Seen" / "\Draft" / flag-keyword / flag-extension
				                    ; Does not include "\Recent"
				flag-extension  = "\" atom
				                   ; Future expansion.  Client implementations
				                   ; MUST accept flag-extension flags.  Server
				                   ; implementations MUST NOT generate
				                   ; flag-extension flags except as defined by
				                   ; future standard or standards-track
				                   ; revisions of this specification.
				flag-keyword    = atom
				date-time       = DQUOTE date-day-fixed "-" date-month "-" date-year SP time SP zone DQUOTE
			*/
			if len(lexCommand.Arguments) < 2 {
				err = errors.New("Parser: expected at least two arguments for APPEND command")
				return
			}
			if !isMailbox(lexCommand.Arguments[0]) {
				err = errors.New("Parser: expected first argument (mailbox) for APPEND to be 'INBOX' or astring")
				return
			}
			if !isLiteral(lexCommand.Arguments[len(lexCommand.Arguments)-1]) {
				err = errors.New("Parser: expected last argument for APPEND to be literal")
				return
			}

			flags := []string{}
			date := time.Time{}

			if len(lexCommand.Arguments) > 2 {
				count := 1
				if lexCommand.Arguments[1][0] == '(' {
					// flag-list
					for i, flag := range lexCommand.Arguments[1 : len(lexCommand.Arguments)-1] {
						flags = append(flags, strings.TrimPrefix(strings.TrimSuffix(flag, ")"), "("))
						if flag[len(flag)-1] == ')' {
							// end
							count = i + 2
							break
						}
					}
					if count == 1 {
						err = errors.New("Parser: malformed flaglist APPEND")
					}
				}
				if count != len(lexCommand.Arguments)-1 {
					// date-time is split into args by lexer
					dateTime := strings.Join(lexCommand.Arguments[count:len(lexCommand.Arguments)-1], " ")
					if !isDateTime(dateTime) {
						err = errors.New("Parser: invalid date-time argument for APPEND")
						return
					}
					dateTime = dateTime[1 : len(dateTime)-1]
					date, err = parseDateTime(strings.TrimPrefix(dateTime, " "))
				}
			}

			command = AppendCmd{
				Mailbox:  parseMailbox(lexCommand.Arguments[0]),
				Literal:  lexCommand.Arguments[len(lexCommand.Arguments)-1],
				Flags:    flags,
				DateTime: date,
			}

		}

	default:
		{
			err = errors.New("Unknown command")
		}
	}

	return

}

func parseMailbox(s string) string {
	if strings.ToUpper(s) == "INBOX" {
		return "INBOX"
	} else {
		return s
	}
}

func parseDateTime(s string) (time.Time, error) {
	return time.Parse("2-Jan-2006 15:04:05 -0700", s)
}
