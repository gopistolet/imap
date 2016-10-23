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

		Convey("Any State", func() {

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

		Convey("Not Authenticated State", func() {

			Convey("STARTTLS", func() {

				cmd, err := parseLine("a001 STARTTLS")
				So(err, ShouldEqual, nil)
				So(cmd, ShouldHaveSameTypeAs, StarttlsCmd{})

				cmd, err = parseLine("a001 STARTTLS no arguments expected")
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

				// Wrong type of arguments
				cmd, err = parseLine("a001 login {test test")
				So(err, ShouldNotEqual, nil)

				cmd, err = parseLine("a001 login test \"test")
				So(err, ShouldNotEqual, nil)

			})

			Convey("AUTHENTICATE", func() {

				cmd, err := parseLine("a001 AUTHENTICATE GSSAPI")
				So(err, ShouldEqual, nil)
				So(cmd, ShouldHaveSameTypeAs, AuthenticateCmd{})
				authCmd := cmd.(AuthenticateCmd)
				So(authCmd.Mechanism, ShouldEqual, "GSSAPI")

				// Not enough arguments
				cmd, err = parseLine("a001 AUTHENTICATE")
				So(err, ShouldNotEqual, nil)

				// Too many arguments
				cmd, err = parseLine("a001 AUTHENTICATE to many args")
				So(err, ShouldNotEqual, nil)

				// Non atom argument
				cmd, err = parseLine("a001 AUTHENTICATE ßlaßlaßla")
				So(err, ShouldNotEqual, nil)

			})

		})

		Convey("Authenticated State", func() {

			for _, test := range []struct {
				instance interface{}
				name     string
			}{
				{
					instance: SelectCmd{},
					name:     "SELECT",
				},
				{
					instance: ExamineCmd{},
					name:     "EXAMINE",
				},
				{
					instance: CreateCmd{},
					name:     "CREATE",
				},
				{
					instance: DeleteCmd{},
					name:     "DELETE",
				},
				{
					instance: SubscribeCmd{},
					name:     "SUBSCRIBE",
				},
				{
					instance: UnsubscribeCmd{},
					name:     "UNSUBSCRIBE",
				},
			} {

				Convey(test.name, func() {

					cmd, err := parseLine("a001 " + test.name + " inbox")
					So(err, ShouldEqual, nil)
					So(cmd, ShouldHaveSameTypeAs, test.instance)
					cmd1 := cmd.(AuthenticatedStateCmd)
					So(cmd1.GetMailbox(), ShouldEqual, "INBOX")

					cmd, err = parseLine("a001 " + test.name + " some_inbox")
					So(err, ShouldEqual, nil)
					So(cmd, ShouldHaveSameTypeAs, test.instance)
					cmd1 = cmd.(AuthenticatedStateCmd)
					So(cmd1.GetMailbox(), ShouldEqual, "some_inbox")

					// Not enough arguments
					cmd, err = parseLine("a001 " + test.name + "")
					So(err, ShouldNotEqual, nil)

					// Too many arguments
					cmd, err = parseLine("a001 " + test.name + " to many args")
					So(err, ShouldNotEqual, nil)

					// Non astring argument
					cmd, err = parseLine("a001 " + test.name + " test\"test")
					So(err, ShouldNotEqual, nil)

				})
			}

			Convey("RENAME", func() {

				cmd, err := parseLine("a001 RENAME source_mailbox dest_mailbox")
				So(err, ShouldEqual, nil)
				So(cmd, ShouldHaveSameTypeAs, RenameCmd{})
				cmd1 := cmd.(RenameCmd)
				So(cmd1.SourceMailbox, ShouldEqual, "source_mailbox")
				So(cmd1.DestinationMailbox, ShouldEqual, "dest_mailbox")

				cmd, err = parseLine("a001 RENAME inbox some_mailbox")
				So(err, ShouldEqual, nil)
				So(cmd, ShouldHaveSameTypeAs, RenameCmd{})
				cmd1 = cmd.(RenameCmd)
				So(cmd1.SourceMailbox, ShouldEqual, "INBOX")

				// Not enough arguments
				cmd, err = parseLine("a001 RENAME")
				So(err, ShouldNotEqual, nil)

				// Too many arguments
				cmd, err = parseLine("a001 RENAME to many args")
				So(err, ShouldNotEqual, nil)

				// Non astring argument
				cmd, err = parseLine("a001 RENAME test\"test test")
				So(err, ShouldNotEqual, nil)

				// Non list-mailbox argument
				cmd, err = parseLine("a001 RENAME test test\"test")
				So(err, ShouldNotEqual, nil)

			})

			Convey("LIST", func() {

				cmd, err := parseLine("a001 LIST some_reference some_mailbox")
				So(err, ShouldEqual, nil)
				So(cmd, ShouldHaveSameTypeAs, ListCmd{})
				cmd1 := cmd.(ListCmd)
				So(cmd1.Reference, ShouldEqual, "some_reference")
				So(cmd1.Mailbox, ShouldEqual, "some_mailbox")

				cmd, err = parseLine("a001 LIST inbox some_mailbox")
				So(err, ShouldEqual, nil)
				So(cmd, ShouldHaveSameTypeAs, ListCmd{})
				cmd1 = cmd.(ListCmd)
				So(cmd1.Reference, ShouldEqual, "INBOX")

				// Not enough arguments
				cmd, err = parseLine("a001 LIST")
				So(err, ShouldNotEqual, nil)

				// Too many arguments
				cmd, err = parseLine("a001 LIST to many args")
				So(err, ShouldNotEqual, nil)

				// Non astring argument
				cmd, err = parseLine("a001 LIST test\"test test")
				So(err, ShouldNotEqual, nil)

				// Non list-mailbox argument
				cmd, err = parseLine("a001 LIST test test\"test")
				So(err, ShouldNotEqual, nil)

			})

			Convey("LSUB", func() {

				cmd, err := parseLine("a001 LSUB some_reference some_mailbox")
				So(err, ShouldEqual, nil)
				So(cmd, ShouldHaveSameTypeAs, LsubCmd{})
				cmd1 := cmd.(LsubCmd)
				So(cmd1.Reference, ShouldEqual, "some_reference")
				So(cmd1.Mailbox, ShouldEqual, "some_mailbox")

				cmd, err = parseLine("a001 LSUB inbox some_mailbox")
				So(err, ShouldEqual, nil)
				So(cmd, ShouldHaveSameTypeAs, LsubCmd{})
				cmd1 = cmd.(LsubCmd)
				So(cmd1.Reference, ShouldEqual, "INBOX")

				// Not enough arguments
				cmd, err = parseLine("a001 LSUB")
				So(err, ShouldNotEqual, nil)

				// Too many arguments
				cmd, err = parseLine("a001 LSUB to many args")
				So(err, ShouldNotEqual, nil)

				// Non astring argument
				cmd, err = parseLine("a001 LSUB test\"test test")
				So(err, ShouldNotEqual, nil)

				// Non list-mailbox argument
				cmd, err = parseLine("a001 LSUB test test\"test")
				So(err, ShouldNotEqual, nil)

			})

		})

	})

}
