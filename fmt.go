package terminfo

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

var (
	errCapEscapeSeq = errors.New("illegal escape sequence")
)

var (
	dvars [26]interface{}
	svars [26]interface{}
)

func (ti *Terminfo) Fmt(cap int, a ...interface{}) ([]byte, error) {
	b, ok := ti.GetStr(cap)
	if !ok {
		return nil, ErrCapAbsent
	}

	b, err := Fmt(b, a...)
	if err == io.EOF {
		err = errCapEscapeSeq
	}
	return b, err
}

// https://invisible-island.net/ncurses/man/terminfo.5.html#h3-Parameterized-Strings
func Fmt(format []byte, params ...interface{}) ([]byte, error) {
	in := bytes.NewReader(format)
	out := bytes.NewBuffer(nil)
	stack := newStack()

	for {
		b, err := in.ReadByte()
		if err == io.EOF {
			if !stack.isEmpty() {
				return nil, errStackNotEmpty
			}
			return out.Bytes(), nil
		}

		if err != nil {
			return nil, err
		}

		if b != '%' {
			out.WriteByte(b)
			continue
		}

		// % codes

		b, err = in.ReadByte()
		if err != nil {
			return nil, err
		}

		switch b {
		case '%':
			out.WriteByte(b)

		case 'c':
			n, err := stack.popI()
			if err != nil {
				return nil, err
			}
			out.WriteString(fmt.Sprintf("%c", n))

		case 'p':
			b, err = in.ReadByte()
			if err != nil {
				return nil, err
			}
			if b < '1' || b > '9' {
				return nil, errors.New("illegal escape sequence: %p[1-9]")
			}
			i := b - '1'
			if int(i) >= len(params) {
				return nil, errors.New("illegal escape sequence: %p[1-9], illegal param idx")
			}
			stack.push(params[i])

		case 'P':
			b, err = in.ReadByte()
			if err != nil {
				return nil, err
			}
			e, err := stack.pop()
			if err != nil {
				return nil, err
			}

			if b >= 'a' && b <= 'z' {
				dvars[b-'a'] = e
			} else if b >= 'A' && b <= 'Z' {
				svars[b-'A'] = e
			}
			return nil, errors.New("illegal escape sequence: %P[a-z|A-Z]")

		case 'g':
			b, err = in.ReadByte()
			if err != nil {
				return nil, err
			}
			if b >= 'a' && b <= 'z' {
				stack.push(dvars[b-'a'])
			} else if b >= 'A' && b <= 'Z' {
				stack.push(svars[b-'A'])
			}
			return nil, errors.New("illegal escape sequence: %g[a-z|A-Z]")

		case '\'':
			b, err = in.ReadByte()
			if err != nil {
				return nil, err
			}
			q, err := in.ReadByte()
			if err != nil {
				return nil, err
			}
			if q != '\'' {
				return nil, errors.New("illegal escape sequence: %'c'")
			}
			stack.push(int(b))

		case '{':
			n := 0
			for {
				b, err = in.ReadByte()
				if err != nil {
					return nil, err
				}

				if b == '}' {
					break
				}

				if b < '0' || b > '9' {
					return nil, errors.New("illegal escape sequence: %{nn}")
				}

				n *= 10
				n += int(b - '0')
			}
			stack.push(n)

		case 'l':
			s, err := stack.popS()
			if err != nil {
				return nil, err
			}
			stack.push(len(s))

		case '+':
			n1, n2, err := stack.popII()
			if err != nil {
				return nil, err
			}
			stack.push(n1 + n2)

		case '-':
			n1, n2, err := stack.popII()
			if err != nil {
				return nil, err
			}
			stack.push(n2 - n1)

		case '*':
			n1, n2, err := stack.popII()
			if err != nil {
				return nil, err
			}
			stack.push(n1 * n2)

		case '/':
			n1, n2, err := stack.popII()
			if err != nil {
				return nil, err
			}
			stack.push(n2 / n1)

		case 'm':
			n1, n2, err := stack.popII()
			if err != nil {
				return nil, err
			}
			stack.push(n2 % n1)

		case '&':
			n1, n2, err := stack.popII()
			if err != nil {
				return nil, err
			}
			stack.push(n1 & n2)

		case '|':
			n1, n2, err := stack.popII()
			if err != nil {
				return nil, err
			}
			stack.push(n1 | n2)

		case '^':
			n1, n2, err := stack.popII()
			if err != nil {
				return nil, err
			}
			stack.push(n1 ^ n2)

		case '=':
			n1, n2, err := stack.popII()
			if err != nil {
				return nil, err
			}
			if n1 == n2 {
				stack.push(1)
			} else {
				stack.push(0)
			}

		case '<':
			n1, n2, err := stack.popII()
			if err != nil {
				return nil, err
			}
			if n2 < n1 {
				stack.push(1)
			} else {
				stack.push(0)
			}

		case '>':
			n1, n2, err := stack.popII()
			if err != nil {
				return nil, err
			}
			if n2 > n1 {
				stack.push(1)
			} else {
				stack.push(0)
			}

		case 'A':
			n1, n2, err := stack.popII()
			if err != nil {
				return nil, err
			}
			if n1 != 0 && n2 != 0 {
				stack.push(1)
			} else {
				stack.push(0)
			}

		case 'O':
			n1, n2, err := stack.popII()
			if err != nil {
				return nil, err
			}
			if n1 != 0 || n2 != 0 {
				stack.push(1)
			} else {
				stack.push(0)
			}

		case '!':
			n, err := stack.popI()
			if err != nil {
				return nil, err
			}

			if n == 0 {
				stack.push(1)
			} else {
				stack.push(0)
			}

		case '~':
			n, err := stack.popI()
			if err != nil {
				return nil, err
			}
			stack.push(^n)

		case 'i':
			if len(params) < 2 {
				return nil, errors.New("illegal escape sequence: %i")
			}

			n1, ok := params[0].(int)
			if !ok {
				return nil, errors.New("illegal escape sequence: %i")
			}

			n2, ok := params[1].(int)
			if !ok {
				return nil, errors.New("illegal escape sequence: %i")
			}

			params[0] = n1 + 1
			params[1] = n2 + 1

		case '?':
			// nop
		case 't':
			n, err := stack.popI()
			if err != nil {
				return nil, err
			}

			if n != 0 {
				break // thenpart
			}

			// goto elsepart

			level := 0
			for {
				b, err = in.ReadByte()
				if err != nil {
					return nil, err
				}

				if b != '%' {
					continue
				}

				b, err = in.ReadByte()
				if err != nil {
					return nil, err
				}

				if b == '?' {
					level++
					continue
				}

				if b == ';' {
					if level > 0 {
						level--
						continue
					}
					break
				}

				if b == 'e' && level == 0 {
					break
				}
			}

		case 'e':
			// skip elsepart
			level := 0
			for {
				b, err = in.ReadByte()
				if err != nil {
					return nil, err
				}

				if b != '%' {
					continue
				}

				b, err = in.ReadByte()
				if err != nil {
					return nil, err
				}

				if b == '?' {
					level++
					continue
				}

				if b == ';' {
					if level > 0 {
						level--
						continue
					}
					break
				}
			}

		case ';':
			// nop

		default:
			f, b, err := parseVerb(b, in)
			if err != nil {
				return nil, err
			}

			if b == 's' {
				s, err := stack.popS()
				if err != nil {
					return nil, err
				}
				out.WriteString(fmt.Sprintf(f, s))
			} else {
				n, err := stack.popI()
				if err != nil {
					return nil, err
				}
				out.WriteString(fmt.Sprintf(f, n))
			}
		}
	}
}

func parseVerb(b byte, in io.ByteReader) (string, byte, error) {
	buf := bytes.NewBuffer(nil)
	buf.WriteByte('%')

	var err error
	for {
		switch b {
		case ':':
			if buf.Len() > 0 {
				return "", 0, errors.New("illegal escape sequence: verb")
			}
		case '-', '+', '#', ' ':
			fallthrough

		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			fallthrough

		case '.':
			fallthrough

		case 'd', 'o', 'x', 'X', 's':
			buf.WriteByte(b)
			return buf.String(), b, nil

		default:
			return "", 0, errors.New("illegal format")
		}

		b, err = in.ReadByte()
		if err != nil {
			return "", 0, err
		}
	}
}
