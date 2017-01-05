package Utils

type RETCODE int

const (
	UNKWON    RETCODE = -2 // UNKWON == -2 (untyped integer constant)
	FAIL      RETCODE = -1 // FAIL == -1 (untyped integer constant)
	SUCC      RETCODE = 0  // SUCC == 0 (untyped integer constant)
	SERVERERR RETCODE = 1  // SERVERERR == 1 (untyped integer constant)
)

type SERVERMODE int

const (
	TEST    SERVERMODE = 0 // TEST == 0 (untyped integer constant)
	DEV     SERVERMODE = 1 // DEV == 1 (untyped integer constant)
	PRODUCT SERVERMODE = 2 // PRODUCT == 2 (untyped integer constant)
)
