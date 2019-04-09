package terminfo

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func LoadEnv() (*Terminfo, error) {
	term := os.Getenv("TERM")
	if term == "" {
		return nil, errors.New("env TERM is empty")
	}
	return Load(term)
}

func Load(term string) (*Terminfo, error) {
	bs, fp, err := LoadTermFile(term)
	if err != nil {
		return nil, err
	}
	return Parse(bs, term, fp)
}

// https://invisible-island.net/ncurses/man/terminfo.5.html#h3-Fetching-Compiled-Descriptions
func LoadTermFile(term string) ([]byte, string, error) {
	if term == "" {
		return nil, "", errors.New("illegal term name")
	}

	dir := os.Getenv("TERMINFO")
	if dir != "" {
		bs, fp, err := ReadTermFile(dir, term)
		if err == nil {
			return bs, fp, nil
		}
	}

	dir = os.Getenv("HOME")
	if dir != "" {
		dir += "/.terminfo"
		bs, fp, err := ReadTermFile(dir, term)
		if err == nil {
			return bs, fp, nil
		}
	}

	dirS := os.Getenv("TERMINFO_DIRS")
	if dirS != "" {
		dirs := strings.Split(dirS, ":")
		for i := 0; i < len(dirs); i++ {
			dir = dirs[i]
			if dir == "" {
				dir = "/usr/share/terminfo"
			}
			bs, fp, err := ReadTermFile(dir, term)
			if err == nil {
				return bs, fp, nil
			}
		}
	}

	bs, fp, err := ReadTermFile("/usr/local/ncurses/share/terminfo", term)
	if err == nil {
		return bs, fp, nil
	}

	bs, fp, err = ReadTermFile("/usr/share/terminfo", term)
	if err == nil {
		return bs, fp, nil
	}

	return nil, "", errors.New("term not found")
}

func ReadTermFile(dir, term string) ([]byte, string, error) {
	fp := fmt.Sprintf("%s/%c/%s", dir, []rune(term)[0], term)
	fp = path.Clean(fp)
	bs, err := ioutil.ReadFile(fp)
	if err != nil {
		return nil, "", err
	}
	return bs, fp, nil
}
