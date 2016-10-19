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

}
