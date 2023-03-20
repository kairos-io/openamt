package amtrpc

type Interface interface {
	CheckAccess() int
	Exec(command string) (string, int)
}
