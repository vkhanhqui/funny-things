package log

import (
	"os"

	olog "log"
)

var (
	StdOutLog = olog.New(os.Stdout, "", 0)
	StdErrLog = olog.New(os.Stderr, "", 0)
	NullLog   = olog.New(NewNullWriter(), "", 0)
	NoLog     = olog.New(NewNullWriter(), "", 0)
	NoLogger  = olog.New(NewNullWriter(), "", 0)
)
