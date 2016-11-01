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

			Convey("STATUS", func() {

				cmd, err := parseLine("A042 STATUS blurdybloop (UIDNEXT MESSAGES)")
				So(err, ShouldEqual, nil)
				So(cmd, ShouldHaveSameTypeAs, StatusCmd{})
				cmd1 := cmd.(StatusCmd)
				So(cmd1.Mailbox, ShouldEqual, "blurdybloop")
				So(len(cmd1.StatusAttributes), ShouldEqual, 2)
				So(cmd1.StatusAttributes[0], ShouldEqual, "UIDNEXT")
				So(cmd1.StatusAttributes[1], ShouldEqual, "MESSAGES")

				cmd, err = parseLine("A042 STATUS blurdybloop (UNSEEN)")
				So(err, ShouldEqual, nil)
				So(cmd, ShouldHaveSameTypeAs, StatusCmd{})
				cmd1 = cmd.(StatusCmd)
				So(cmd1.Mailbox, ShouldEqual, "blurdybloop")
				So(len(cmd1.StatusAttributes), ShouldEqual, 1)
				So(cmd1.StatusAttributes[0], ShouldEqual, "UNSEEN")

				// Not enough arguments
				cmd, err = parseLine("a001 STATUS")
				So(err, ShouldNotEqual, nil)

				cmd, err = parseLine("a001 STATUS box")
				So(err, ShouldNotEqual, nil)

				// Too many arguments
				cmd, err = parseLine("a001 STATUS to many args")
				So(err, ShouldNotEqual, nil)

				// Non astring argument
				cmd, err = parseLine("a001 STATUS test\"test test")
				So(err, ShouldNotEqual, nil)

				// Non attr list argument
				cmd, err = parseLine("a001 STATUS test test")
				So(err, ShouldNotEqual, nil)

				// Empty attr list argument
				cmd, err = parseLine("a001 STATUS test ()")
				So(err, ShouldNotEqual, nil)

				// Wrong attr argument
				cmd, err = parseLine("a001 STATUS test (somename)")
				So(err, ShouldNotEqual, nil)

				// wrong syntax for attr argument
				cmd, err = parseLine("a001 STATUS test (somename")
				So(err, ShouldNotEqual, nil)

				cmd, err = parseLine("a001 STATUS test somename)")
				So(err, ShouldNotEqual, nil)

			})

			Convey("APPEND", func() {

				cmd, err := parseLine("A003 APPEND saved-messages (\\Seen) {310}")
				So(err, ShouldEqual, nil)
				So(cmd, ShouldHaveSameTypeAs, AppendCmd{})
				cmd1 := cmd.(AppendCmd)
				So(cmd1.Mailbox, ShouldEqual, "saved-messages")
				So(cmd1.Flags, ShouldResemble, []string{"\\Seen"})

				cmd, err = parseLine(`A00027 APPEND A-SPAM-filtered/2002 (\Seen) "31-Dec-2002 14:36:36 -0800" {6663}`)
				So(err, ShouldEqual, nil)

				cmd, err = parseLine(`A00027 APPEND A-SPAM-filtered/2002 "31-Dec-2002 14:36:36 -0800" {6663}`)
				So(err, ShouldEqual, nil)
				So(cmd, ShouldHaveSameTypeAs, AppendCmd{})
				cmd1 = cmd.(AppendCmd)
				So(cmd1.Mailbox, ShouldEqual, "A-SPAM-filtered/2002")
				So(cmd1.Flags, ShouldResemble, []string{})

				cmd, err = parseLine(`A00027 APPEND A-SPAM-filtered/2002 " 1-Dec-2002 14:36:36 +0800" {6663}`)
				So(err, ShouldEqual, nil)
				So(cmd, ShouldHaveSameTypeAs, AppendCmd{})
				cmd1 = cmd.(AppendCmd)
				So(cmd1.Mailbox, ShouldEqual, "A-SPAM-filtered/2002")
				So(cmd1.Flags, ShouldResemble, []string{})

				// Not enough arguments
				cmd, err = parseLine("a001 APPEND")
				So(err, ShouldNotEqual, nil)

				cmd, err = parseLine("a001 APPEND box")
				So(err, ShouldNotEqual, nil)

				// Too many arguments
				cmd, err = parseLine("a001 APPEND to many args")
				So(err, ShouldNotEqual, nil)

				// Non mailbox argument
				cmd, err = parseLine("a001 APPEND test\"test {10}")
				So(err, ShouldNotEqual, nil)

				// Non literal argument
				cmd, err = parseLine("a001 APPEND test {test")
				So(err, ShouldNotEqual, nil)

				// malformed list
				cmd, err = parseLine("A003 APPEND saved-messages (\\Seen {310}")
				So(err, ShouldNotEqual, nil)

				// malformed time
				cmd, err = parseLine("A003 APPEND saved-messages (\\Seen) \"1-malformed-2002 14:36:36 +0800\" {310}")
				So(err, ShouldNotEqual, nil)
			})

		})

		Convey("Selected State", func() {

			Convey("CHECK", func() {

				cmd, err := parseLine("FXXZ CHECK")
				So(err, ShouldEqual, nil)
				So(cmd, ShouldHaveSameTypeAs, CheckCmd{})

				cmd, err = parseLine("a001 CHECK no arguments expected")
				So(err, ShouldNotEqual, nil)
			})

			Convey("CLOSE", func() {

				cmd, err := parseLine("A341 CLOSE")
				So(err, ShouldEqual, nil)
				So(cmd, ShouldHaveSameTypeAs, CloseCmd{})

				cmd, err = parseLine("a001 CLOSE no arguments expected")
				So(err, ShouldNotEqual, nil)
			})

			Convey("EXPUNGE", func() {

				cmd, err := parseLine("A202 EXPUNGE")
				So(err, ShouldEqual, nil)
				So(cmd, ShouldHaveSameTypeAs, ExpungeCmd{})

				cmd, err = parseLine("a001 EXPUNGE no arguments expected")
				So(err, ShouldNotEqual, nil)
			})

			Convey("FETCH", func() {

				cmd, err := parseLine("A654 FETCH 2:4 (FLAGS BODY[HEADER.FIELDS (DATE FROM)])")
				So(err, ShouldEqual, nil)
				So(cmd, ShouldHaveSameTypeAs, FetchCmd{})
			})

			Convey("Store", func() {

				cmd, err := parseLine("A003 STORE 2:4 +FLAGS (\\Deleted)")
				So(err, ShouldEqual, nil)
				So(cmd, ShouldHaveSameTypeAs, StoreCmd{})
				cmd1 := cmd.(StoreCmd)
				So(cmd1.Sequence, ShouldEqual, "2:4")
				So(cmd1.Mode, ShouldEqual, "+")
				So(cmd1.Silent, ShouldEqual, false)
				So(cmd1.Flags, ShouldResemble, []string{"\\Deleted"})

				cmd, err = parseLine("A003 STORE 2:4 FLAGS.SILENT (\\Deleted \\Seen)")
				So(err, ShouldEqual, nil)
				So(cmd, ShouldHaveSameTypeAs, StoreCmd{})
				cmd1 = cmd.(StoreCmd)
				So(cmd1.Sequence, ShouldEqual, "2:4")
				So(cmd1.Mode, ShouldEqual, "")
				So(cmd1.Silent, ShouldEqual, true)
				So(cmd1.Flags, ShouldResemble, []string{"\\Deleted", "\\Seen"})

				// Not enough args
				cmd, err = parseLine("A003 STORE 2:4")
				So(err, ShouldNotEqual, nil)

				cmd, err = parseLine("A003 STORE")
				So(err, ShouldNotEqual, nil)

				// not sequence set as first arg
				cmd, err = parseLine("A003 STORE blablabla FLAGS.SILENT (\\Deleted \\Seen)")
				So(err, ShouldNotEqual, nil)

				// No flags
				cmd, err = parseLine("A003 STORE 2:4 somekey.SILENT (\\Deleted \\Seen)")
				So(err, ShouldNotEqual, nil)

				// wrong mode
				cmd, err = parseLine("A003 STORE 2:4 mFLAGS (\\Deleted)")
				So(err, ShouldNotEqual, nil)

				// malformed list
				cmd, err = parseLine("A003 STORE 2:4 FLAGS \\Deleted)")
				So(err, ShouldNotEqual, nil)

				cmd, err = parseLine("A003 STORE 2:4 FLAGS (\\Deleted")
				So(err, ShouldNotEqual, nil)
			})

		})

	})

}
