package parser

import (
	"errors"
	"strings"
)

type Argument string

type lexCommand struct {
	Tag       string
	Name      string
	Arguments []string
}

// lexLine creates a command struct for an IMAP line
// which contains the Name of the command and the arguments
func lexLine(line string) (c lexCommand, err error) {
	parts := strings.Split(line, " ")
	if len(parts) >= 2 {
		if !isTag(parts[0]) {
			err = errors.New("Lexer: expected identifier tag")
			return
		}
		c.Tag = parts[0]
		if !isCommand(parts[1]) {
			err = errors.New("Lexer: expected valid IMAP command")
			return
		}
		c.Name = strings.ToUpper(parts[1])
	}
	if len(parts) > 2 {
		c.Arguments = parts[2:]
	}
	return
}

/*
tag = 1*<any ASTRING-CHAR except "+">
*/
func isTag(s string) bool {
	for _, c := range s {
		if !isAStringChar(c) {
			return false
		}
		if c == '+' {
			return false
		}
	}
	return true
}

/*
astring         = 1*ASTRING-CHAR / string
string          = quoted / literal

quoted          = DQUOTE *QUOTED-CHAR DQUOTE

literal         = "{" number "}" CRLF *CHAR8
                  ; Number represents the number of CHAR8s
number          = 1*DIGIT
                  ; Unsigned 32-bit integer
                  ; (0 <= n < 4,294,967,296)
CRLF            =  %d13.10
CHAR8           = %x01-ff
                  ; any OCTET except NUL, %x00
*/
func isAString(s string) bool {
	if len(s) == 0 {
		return false
	}

	switch s[0] {
	case '"':
		{
			if len(s) < 2 {
				return false
			}
			if s[len(s)-1] != '"' {
				return false
			}
			return isQuoted(s[1 : len(s)-1])
		}
	case '{':
		{
			// TODO: CRLF *CHAR8
			if s[len(s)-1] != '}' {
				return false
			}
			for _, c := range s[1 : len(s)-1] {
				if !isDigit(c) {
					return false
				}
			}
		}
	default:
		{
			for _, c := range s {
				if !isAStringChar(c) {
					return false
				}
			}
		}
	}

	return true
}

/*
quoted          = DQUOTE *QUOTED-CHAR DQUOTE

QUOTED-CHAR     = <any TEXT-CHAR except quoted-specials> /
				  "\" quoted-specials
quoted-specials = DQUOTE / "\"
*/
func isQuoted(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' {
			if !(len(s) > i+1) {
				return false
			}
			if s[i+1] != '"' && s[i+1] != '\\' {
				return false
			}
			i++
		} else if s[i] == '"' {
			return false
		}
	}
	return true
}

/*
DIGIT  =  %x30-39
	      ; 0-9
*/
func isDigit(c rune) bool {
	if c >= 0x30 && c <= 0x39 {
		return true
	} else {
		return false
	}
}

func isAlpha(c rune) bool {
	if c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' {
		return true
	} else {
		return false
	}
}

func isCommand(s string) bool {
	for _, c := range s {
		if !isAlpha(c) {
			return false
		}
	}
	return true
}

/*
ASTRING-CHAR    = ATOM-CHAR / resp-specials
resp-specials   = "]"
*/
func isAStringChar(c rune) bool {
	if isAtomChar(c) {
		return true
	}
	if c == ']' {
		return true
	}
	return false
}

/*
atom            = 1*ATOM-CHAR
ATOM-CHAR       = <any CHAR except atom-specials>
atom-specials   = "(" / ")" / "{" / SP / CTL / list-wildcards /
				  quoted-specials / resp-specials
CHAR            =  %x01-7F
                     ; any 7-bit US-ASCII character,
                     ;  excluding NUL
*/
func isAtomChar(c rune) bool {
	switch c {
	case '(', ')', '{', ' ', '%', '*', '"', '\\', ']':
		return false
	case 0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xa, 0xb, 0xc, 0xd, 0xe, 0xf, 0x10, 0x11, 0x12,
		0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x7f: // CTL chars
		return false
	}
	if c >= 0x01 && c <= 0x7f {
		return true
	} else {
		return false
	}
}

func isAtom(s string) bool {
	for _, c := range s {
		if !isAtomChar(c) {
			return false
		}
	}
	return true
}

/*
mailbox = "INBOX" / astring
		  ; INBOX is case-insensitive.  All case variants of
		  ; INBOX (e.g., "iNbOx") MUST be interpreted as INBOX
		  ; not as an astring.  An astring which consists of
		  ; the case-insensitive sequence "I" "N" "B" "O" "X"
		  ; is considered to be INBOX and not an astring.
		  ;  Refer to section 5.1 for further
		  ; semantic details of mailbox names.
*/
func isMailbox(s string) bool {
	if strings.ToUpper(s) == "INBOX" {
		return true
	} else if isAString(s) {
		return true
	} else {
		return false
	}
}

/*
list-mailbox    = 1*list-char / string
list-char       = ATOM-CHAR / list-wildcards / resp-specials
list-wildcards  = "%" / "*"

string          = quoted / literal
*/
func isListMailbox(s string) bool {
	if len(s) == 0 {
		return false
	}

	switch s[0] {
	case '"':
		{
			// string -> quoted
			if len(s) < 2 {
				return false
			}
			if s[len(s)-1] != '"' {
				return false
			}
			return isQuoted(s[1 : len(s)-1])
		}
	case '{':
		{
			// string -> literal
			if s[len(s)-1] != '}' {
				return false
			}
			for _, c := range s[1 : len(s)-1] {
				if !isDigit(c) {
					return false
				}
			}
		}
	default:
		{
			// 1*list-char
			// TODO: CRLF *CHAR8
			for _, c := range s {
				if !isAtomChar(c) && c != '%' && c != '*' && c != ']' {
					return false
				}
			}
		}
	}

	return true
}
