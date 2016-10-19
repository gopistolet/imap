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
tag             = 1*<any ASTRING-CHAR except "+">
ASTRING-CHAR    = ATOM-CHAR / resp-specials
ATOM-CHAR       = <any CHAR except atom-specials>
atom-specials   = "(" / ")" / "{" / SP / CTL / list-wildcards / quoted-specials / resp-specials
list-wildcards  = "%" / "*"
quoted-specials = DQUOTE / "\"
resp-specials   = "]"
CTL             =  %x00-1F / %x7F
*/
func isTag(s string) bool {
	for _, c := range s {
		switch c {
		case '(', ')', '{', ' ', '%', '*', '"', '\\':
			return false
		case 0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xa, 0xb, 0xc, 0xd, 0xe, 0xf, 0x10, 0x11, 0x12,
			0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x7f: // CTL chars
			return false
		}
	}
	return true
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
