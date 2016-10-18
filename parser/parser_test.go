package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestParser(t *testing.T) {

	Convey("Testing parseLine", t, func() {

		Convey("Testing general stuff", func() {

			_, err := parseLine("a001 n00p")
			So(err, ShouldNotEqual, nil)

			_, err = parseLine("a001 unknowncommand")
			So(err, ShouldNotEqual, nil)

		})

		Convey("LOGIN", func() {

			cmd, err := parseLine("a001 login mrc secret")
			So(err, ShouldEqual, nil)
			So(cmd, ShouldHaveSameTypeAs, LoginCmd{})
			loginCmd := cmd.(LoginCmd)
			So(loginCmd.Username, ShouldEqual, "mrc")
			So(loginCmd.Password, ShouldEqual, "secret")

			// Not enough arguments
			cmd, err = parseLine("a001 login")
			So(err, ShouldNotEqual, nil)

			// Not enough arguments
			cmd, err = parseLine("a001 login test")
			So(err, ShouldNotEqual, nil)

		})

		Convey("LOGOUT", func() {

			cmd, err := parseLine("a001 LOGOUT")
			So(err, ShouldEqual, nil)
			So(cmd, ShouldHaveSameTypeAs, LogoutCmd{})

			cmd, err = parseLine("a001 LOGOUT no arguments expected")
			So(err, ShouldNotEqual, nil)
		})

		Convey("CAPABILITY", func() {

			cmd, err := parseLine("a001 CAPABILITY")
			So(err, ShouldEqual, nil)
			So(cmd, ShouldHaveSameTypeAs, CapabilityCmd{})

			cmd, err = parseLine("a001 CAPABILITY no arguments expected")
			So(err, ShouldNotEqual, nil)

		})

		Convey("NOOP", func() {

			cmd, err := parseLine("a001 NOOP")
			So(err, ShouldEqual, nil)
			So(cmd, ShouldHaveSameTypeAs, NoopCmd{})

			cmd, err = parseLine("a001 NOOP no arguments expected")
			So(err, ShouldNotEqual, nil)

		})

	})

}
