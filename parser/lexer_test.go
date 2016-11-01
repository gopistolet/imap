package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestLexer(t *testing.T) {

	Convey("Testing lexLine", t, func() {

		c, err := lexLine("a003 Fetch 12 full")
		So(err, ShouldEqual, nil)
		So(c.Name, ShouldEqual, "FETCH")
		So(c.Tag, ShouldEqual, "a003")
		So(c.Arguments[0], ShouldEqual, "12")
		So(c.Arguments[1], ShouldEqual, "full")
		So(len(c.Arguments), ShouldEqual, 2)

		c, err = lexLine("a002 NOOP")
		So(err, ShouldEqual, nil)
		So(c.Name, ShouldEqual, "NOOP")
		So(c.Tag, ShouldEqual, "a002")
		So(len(c.Arguments), ShouldEqual, 0)

		// Invalid command
		c, err = lexLine("a002 n00b")
		So(err, ShouldNotEqual, nil)

		// Invalid tag
		c, err = lexLine("\\a002 test")
		So(err, ShouldNotEqual, nil)

	})

	Convey("Testing isCommand", t, func() {
		for _, command := range []string{
			"fetch",
			"FETCH",
			"Fetch",
		} {
			So(isCommand(command), ShouldEqual, true)
		}

		for _, command := range []string{
			"a002",
			"test+test",
			"Fetch&",
		} {
			So(isCommand(command), ShouldEqual, false)
		}
	})

	Convey("Testing isTag", t, func() {
		for _, command := range []string{
			"a002",
			"test",
			"1",
			"test]",
		} {
			So(isTag(command), ShouldEqual, true)
		}

		for _, command := range []string{
			`"invalid"`,
			"test*test",
			"test%test",
			"{test}",
			"(test)",
			" ",
			"\\test",
			"test+test",
			string([]byte{0x0, 0x1, 0x2, 0x3, 0x4}),
			string([]byte{0x7f}),
		} {
			So(isTag(command), ShouldEqual, false)
		}
	})

	Convey("Testing isAtom", t, func() {
		for _, command := range []string{
			"a002",
			"test",
			"1",
			"test",
		} {
			So(isAtom(command), ShouldEqual, true)
		}

		for _, command := range []string{
			`"invalid"`,
			"test*test",
			"test%test",
			"{test}",
			"(test)",
			" ",
			"\\test",
			"]",
			"👎",
			"π",
			string([]byte{0x0, 0x1, 0x2, 0x3, 0x4}),
			string([]byte{0x7f}),
		} {
			So(isAtom(command), ShouldEqual, false)
		}
	})

	Convey("Testing isDigit", t, func() {
		for _, d := range []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'} {
			So(isDigit(d), ShouldEqual, true)
		}

		for _, d := range []rune{'a', 'π', '-', ')'} {
			So(isDigit(d), ShouldEqual, false)
		}
	})

	Convey("Testing isQuoted", t, func() {
		for _, s := range []string{
			"a002",
			"test",
			"1",
			"test",
			`\\`,
			`\"`,
		} {
			So(isQuoted(s), ShouldEqual, true)
		}

		for _, s := range []string{
			`a0\02`,
			`te"st`,
			`"`,
			`\`,
		} {
			So(isQuoted(s), ShouldNotEqual, true)
		}
	})

	Convey("Testing isAString", t, func() {
		for _, s := range []string{
			"a002",
			"test",
			"1",
			"test",
			`"hello world"`,
			`"\\"`,
			`"\""`,
			`"hello\"world"`,
			`"hello'world"`,
			"{10}",
			"{1}",
		} {
			So(isAString(s), ShouldEqual, true)
		}

		for _, s := range []string{
			"test*test",
			"test%test",
			"{test}",
			"{10",
			"{",
			"(",
			")",
			"(test",
			"test)",
			"(test)",
			" ",
			"",
			`"`,
			`"test`,
			`test"`,
			"\\test",
			"👎",
			"π",
			string([]byte{0x0, 0x1, 0x2, 0x3, 0x4}),
			string([]byte{0x7f}),
		} {
			So(isAString(s), ShouldEqual, false)
		}

	})

	Convey("Testing isMailbox", t, func() {
		for _, s := range []string{
			"a002",
			"test",
			"1",
			"test",
			`"hello world"`,
			`"\\"`,
			`"\""`,
			`"hello\"world"`,
			`"hello'world"`,
			"{10}",
			"{1}",
			"INBOX",
			"inbox",
			"InBoX",
			"~smith/Mail/",
			"archive/",
			"#news.",
			"~smith/Mail/",
			"archive/",
		} {
			So(isMailbox(s), ShouldEqual, true)
		}

		for _, s := range []string{
			"test*test",
			"\\test",
			"👎",
			"π",
		} {
			So(isMailbox(s), ShouldNotEqual, true)
		}
	})

	Convey("Testing isListMailbox", t, func() {
		for _, s := range []string{
			"a002",
			"test",
			"1",
			"test",
			`"hello world"`,
			`"\\"`,
			`"\""`,
			`"hello\"world"`,
			`"hello'world"`,
			"{10}",
			"{1}",
			"test*test",
			"test%test",
			"foo.*",
			"%",
			"comp.mail.*",
			"/usr/doc/foo",
			"~fred/Mail/*",
		} {
			So(isListMailbox(s), ShouldEqual, true)
		}

		for _, s := range []string{
			"{test}",
			"{10",
			"{}",
			"{",
			"(",
			")",
			"(test",
			"test)",
			"(test)",
			" ",
			"",
			`"`,
			`"test`,
			`test"`,
			"\\test",
			"👎",
			"π",
			string([]byte{0x0, 0x1, 0x2, 0x3, 0x4}),
			string([]byte{0x7f}),
		} {
			So(isListMailbox(s), ShouldEqual, false)
		}

	})

	Convey("Testing isLiteral", t, func() {
		for _, s := range []string{
			"{10}",
			"{1}",
		} {
			So(isLiteral(s), ShouldEqual, true)
		}

		for _, s := range []string{
			"{test}",
			"{}",
			"{1",
			"1}",
		} {
			So(isLiteral(s), ShouldNotEqual, true)
		}
	})

	Convey("Testing isDateTime", t, func() {
		for _, s := range []string{
			`"31-Dec-2002 14:36:36 -0800"`,
			`"31-Dec-2002 14:36:36 +0800"`,
			`" 1-Dec-2002 14:36:36 -0800"`,
		} {
			So(isDateTime(s), ShouldEqual, true)
		}

		for _, s := range []string{
			`"31-Dec-2002 14:36:36 +0800`,
			`"3c-Dec-2002 14:36:36 +0800"`,
			`"31-zzz-2002 14:36:36 -0800"`,
			`"c1-Dec-2002 14:36:36 +0800"`,
			`"Dec-2002 14:36:36 -0800"`,
			`"31-Dec 14:36:36 -0800"`,
			`"31-2002 14:36:36 -0800"`,
			`"31-Dec-2002 -0800"`,
			`"31-Dec-2002 14:36:36"`,
			`"31-Decgezgezgez-2002 14:36:36 -0800"`,
			`" 1-Dec-2002 14=36=36 -0800"`,
			`" 1-Dec-2002 14:36:36 =0800"`,
			`" 1-Dec-2002 143636 -0800"`,
		} {
			So(isDateTime(s), ShouldNotEqual, true)
		}
	})

	Convey("Testing isNzDigit", t, func() {
		for _, command := range []rune{
			'1',
			'2',
			'3',
		} {
			So(isNzDigit(command), ShouldEqual, true)
		}

		for _, command := range []rune{
			'0',
			'a',
			'z',
		} {
			So(isNzDigit(command), ShouldEqual, false)
		}
	})

	Convey("Testing isSeqNumber", t, func() {
		for _, command := range []string{
			"2",
			"22",
			"*",
		} {
			So(isSeqNumber(command), ShouldEqual, true)
		}

		for _, command := range []string{
			"0",
			"01",
			"*b",
		} {
			So(isSeqNumber(command), ShouldEqual, false)
		}
	})

	Convey("Testing isSeqRange", t, func() {
		for _, command := range []string{
			"2:4",
			"2:*",
			"12:*",
		} {
			So(isSeqRange(command), ShouldEqual, true)
		}

		for _, command := range []string{
			"a:b",
			"0",
			"*",
		} {
			So(isSeqRange(command), ShouldEqual, false)
		}
	})

	Convey("Testing isSequenceSet", t, func() {
		for _, command := range []string{
			"2,4:7,9,12:*",
			"443:557",
			"4827313:4828442",
			"2",
			"2:4",
		} {
			So(isSequenceSet(command), ShouldEqual, true)
		}

		for _, command := range []string{
			"",
			"a",
			"a:b",
			"*:b",
		} {
			So(isSequenceSet(command), ShouldEqual, false)
		}
	})

}
