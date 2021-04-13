package postgresstore

const(
	duplicateErrorCode = "23505"
	sqlDbSourceError = "SQL sb source error"
)

type DuplicateSourceErr struct {
	Err error
}

func (e *DuplicateSourceErr )Error() string {
	return e.Err.Error()
}

