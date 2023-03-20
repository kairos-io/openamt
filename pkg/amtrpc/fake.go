//go:build fake

package amtrpc

type AMTRPC struct {
	MockAccessStatus func() int
	MockExec         func(string) (string, int)
}

func (f AMTRPC) CheckAccess() int {
	return f.MockAccessStatus()
}

func (f AMTRPC) Exec(command string) (string, int) {
	return f.MockExec(command)
}
