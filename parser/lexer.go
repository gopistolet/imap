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
			return isLiteral(s)
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

/*
literal         = "{" number "}" CRLF *CHAR8
				  ; Number represents the number of CHAR8s
number          = 1*DIGIT
				  ; Unsigned 32-bit integer
				  ; (0 <= n < 4,294,967,296)
*/
func isLiteral(s string) bool {

	/*
		TODO:
		should we accept '+'?
		found this on internet:
		    A00010 APPEND A-SPAM-filtered/2002 (\Seen) "31-Dec-2002 22:36:36 -0800" {6663+}
	*/

	if len(s) < 3 {
		return false
	}
	if s[0] != '{' || s[len(s)-1] != '}' {
		return false
	}
	for _, c := range s[1 : len(s)-1] {
		if !isDigit(c) {
			return false
		}
	}
	return true
}

/*
date-time       = DQUOTE date-day-fixed "-" date-month "-" date-year SP time SP zone DQUOTE

date-day-fixed  = (SP DIGIT) / 2DIGIT
					; Fixed-format version of date-day

date-month      = "Jan" / "Feb" / "Mar" / "Apr" / "May" / "Jun" /
				  "Jul" / "Aug" / "Sep" / "Oct" / "Nov" / "Dec"

date-year       = 4DIGIT

time            = 2DIGIT ":" 2DIGIT ":" 2DIGIT
                    ; Hours minutes seconds

zone            = ("+" / "-") 4DIGIT
                    ; Signed four-digit value of hhmm representing
                    ; hours and minutes east of Greenwich (that is,
                    ; the amount that the given time differs from
                    ; Universal Time).  Subtracting the timezone
                    ; from the given time will give the UT form.
                    ; The Universal Time zone is "+0000".

Example: "31-Dec-2002 14:36:36 -0800"
*/
func isDateTime(s string) bool {
	if s[0] != '"' || s[len(s)-1] != '"' {
		return false
	}
	s = s[1 : len(s)-1]

	if len(s) != 26 {
		return false
	}

	// Digits
	for _, d := range []int{1, 7, 8, 9, 10, 12, 13, 15, 16, 18, 19, 22, 23, 24, 25} {
		if !isDigit(rune(s[d])) {
			return false
		}
	}
	if !isDigit(rune(s[0])) && s[0] != ' ' {
		return false
	}

	// Months
	switch s[3:6] {
	case "Jan", "Feb", "Mar", "Apr", "May", "Jun",
		"Jul", "Aug", "Sep", "Oct", "Nov", "Dec":
		break
	default:
		return false
	}

	// Hour seperators
	if s[14] != ':' || s[17] != ':' {
		return false
	}

	// Timezone sign
	if s[21] != '+' && s[21] != '-' {
		return false
	}

	return true
}

/*
nz-number       = digit-nz *DIGIT
					; Non-zero unsigned 32-bit integer
					; (0 < n < 4,294,967,296)
*/
func isNzDigit(c rune) bool {
	switch c {
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return true
	default:
		return false
	}
}

/*
seq-number      = nz-number / "*"
                    ; message sequence number (COPY, FETCH, STORE
                    ; commands) or unique identifier (UID COPY,
                    ; UID FETCH, UID STORE commands).
                    ; * represents the largest number in use.  In
                    ; the case of message sequence numbers, it is
                    ; the number of messages in a non-empty mailbox.
                    ; In the case of unique identifiers, it is the
                    ; unique identifier of the last message in the
                    ; mailbox or, if the mailbox is empty, the
                    ; mailbox's current UIDNEXT value.
                    ; The server should respond with a tagged BAD
                    ; response to a command that uses a message
                    ; sequence number greater than the number of
                    ; messages in the selected mailbox.  This
                    ; includes "*" if the selected mailbox is empty.
*/
func isSeqNumber(s string) bool {
	if len(s) == 1 && s[0] == '*' {
		return true
	} else {
		for _, c := range s {
			if !isNzDigit(c) {
				return false
			}
		}
	}
	return true
}

/*
seq-range       = seq-number ":" seq-number
                    ; two seq-number values and all values between
                    ; these two regardless of order.
                    ; Example: 2:4 and 4:2 are equivalent and indicate
                    ; values 2, 3, and 4.
                    ; Example: a unique identifier sequence range of
                    ; 3291:* includes the UID of the last message in
                    ; the mailbox, even if that value is less than 3291.
*/
func isSeqRange(s string) bool {
	sp := strings.Split(s, ":")
	if len(sp) != 2 {
		return false
	}
	return isSeqNumber(sp[0]) && isSeqNumber(sp[1])
}

/*
sequence-set    = (seq-number / seq-range) *("," sequence-set)
					; set of seq-number values, regardless of order.
					; Servers MAY coalesce overlaps and/or execute the
					; sequence in any order.
					; Example: a message sequence number set of
					; 2,4:7,9,12:* for a mailbox with 15 messages is
					; equivalent to 2,4,5,6,7,9,12,13,14,15
					; Example: a message sequence number set of *:4,5:7
					; for a mailbox with 10 messages is equivalent to
					; 10,9,8,7,6,5,4,5,6,7 and MAY be reordered and
					; overlap coalesced to be 4,5,6,7,8,9,10.
*/
func isSequenceSet(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, seq := range strings.Split(s, ",") {
		if !isSeqRange(seq) && !isSeqNumber(seq) {
			return false
		}
	}
	return true
}
