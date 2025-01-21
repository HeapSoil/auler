package errs

var (
	ErrSpellNotFound = &Errno{
		HTTP:    404,
		Code:    "ResourceNotFound.SpellNotFound",
		Message: "Spell already exist.",
	}
)
