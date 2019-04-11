package terminfo

import (
	"os"

	"golang.org/x/sys/unix"
)

func Isatty(fd uintptr) bool {
	_, err := unix.IoctlGetTermios(int(fd), unix.TCGETS)
	return err == nil
}

func TermFileIn() (*os.File, error) {
	if Isatty(os.Stdin.Fd()) {
		return os.Stdin, nil
	}

	f, err := os.Open("/dev/tty")
	if err != nil {
		return nil, err
	}

	if Isatty(f.Fd()) {
		return f, nil
	}

	f.Close()
	return nil, os.ErrNotExist
}

func TermFileOut() (*os.File, error) {
	if Isatty(os.Stderr.Fd()) {
		return os.Stderr, nil
	}

	if Isatty(os.Stdout.Fd()) {
		return os.Stdout, nil
	}

	f, err := os.OpenFile("/dev/tty", os.O_WRONLY, 0)
	if err != nil {
		return nil, err
	}

	if Isatty(f.Fd()) {
		return f, nil
	}

	f.Close()
	return nil, os.ErrNotExist
}

func TermFileClose(f *os.File) error {
	if f == os.Stdin || f == os.Stdout || f == os.Stderr {
		return nil
	}
	return f.Close()
}

func TermSize() (int, int, error) { // TODO optimize
	f, err := TermFileOut()
	if err != nil {
		return 0, 0, err
	}
	defer TermFileClose(f)

	ws, err := unix.IoctlGetWinsize(int(f.Fd()), unix.TIOCGWINSZ)
	if err != nil {
		return 0, 0, err
	}
	return int(ws.Row), int(ws.Col), nil
}
