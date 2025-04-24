package ae

var (
	Separator       = " - "
	BadParamFormat  = "Bad Parameter: %s"
	ParamIsTooShort = "Too Short Parameter: %s"
	ParamIsTooLong  = "Too Long Parameter: %s"
	WrongPassword   = "Wrong Password"
	WrongToken      = "Wrong Token"
	TokenExpired    = "Token Expired"
)

var (
	newCodeTexts = map[int]string{
		FailedAndSeeOther:    "Failed And See Other",   // 391
		RequestTargetInvalid: "Request Target Invalid", // 414
		PageExpired:          "Page Expired",           // 419
		EnhanceYourCalm:      "Enhance Your Calm",      // 420
		NoRowsAvailable:      "No Rows Available",      //494
		ConflictWith:         "Conflict With",          //499
		Exception:            "Exception",              // 590
	}
)
