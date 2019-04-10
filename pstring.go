package terminfo

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

// https://invisible-island.net/ncurses/man/terminfo.5.html#h3-Parameterized-Strings
func Fmt(format []byte, params ...interface{}) ([]byte, error) {
	in := bytes.NewReader(format)
	out := new(bytes.Buffer)

	stack := newStack()
	var dvars [26]interface{}
	var svars [26]interface{}

	for {
		b, err := in.ReadByte()
		if err == io.EOF {
			if !stack.isEmpty() {
				return nil, errors.New("stack not empty") // TODO ignore it?
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
			n, err := stack.popN()
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
				return nil, errors.New("%p[1-9]")
			}
			i := b - '1'
			if int(i) >= len(params) {
				return nil, errors.New("%p[1-9] illegal param")
			}
			stack.push(params[i])

		case 'P':
			b, err = in.ReadByte()
			if err != nil {
				return nil, err
			}
			v, err := stack.pop()
			if err != nil {
				return nil, err
			}

			if b >= 'a' && b <= 'z' {
				dvars[b-'a'] = v
			} else if b >= 'A' && b <= 'Z' {
				svars[b-'A'] = v
			}
			return nil, errors.New("%P[a-z] / %P[A-Z]")

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
			return nil, errors.New("%g[a-z] / %g[A-Z]")

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
				return nil, errors.New("%'c'")
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
					return nil, errors.New("%{nn}")
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
			n1, n2, err := stack.popN2()
			if err != nil {
				return nil, err
			}
			stack.push(n1 + n2)

		case '-':
			n1, n2, err := stack.popN2()
			if err != nil {
				return nil, err
			}
			stack.push(n2 - n1)

		case '*':
			n1, n2, err := stack.popN2()
			if err != nil {
				return nil, err
			}
			stack.push(n1 * n2)

		case '/':
			n1, n2, err := stack.popN2()
			if err != nil {
				return nil, err
			}
			stack.push(n2 / n1)

		case 'm':
			n1, n2, err := stack.popN2()
			if err != nil {
				return nil, err
			}
			stack.push(n2 % n1)

		case '&':
			n1, n2, err := stack.popN2()
			if err != nil {
				return nil, err
			}
			stack.push(n1 & n2)

		case '|':
			n1, n2, err := stack.popN2()
			if err != nil {
				return nil, err
			}
			stack.push(n1 | n2)

		case '^':
			n1, n2, err := stack.popN2()
			if err != nil {
				return nil, err
			}
			stack.push(n1 ^ n2)

		case '=':
			n1, n2, err := stack.popN2()
			if err != nil {
				return nil, err
			}
			if n1 == n2 {
				stack.push(1)
			} else {
				stack.push(0)
			}

		case '<':
			n1, n2, err := stack.popN2()
			if err != nil {
				return nil, err
			}
			if n2 < n1 {
				stack.push(1)
			} else {
				stack.push(0)
			}

		case '>':
			n1, n2, err := stack.popN2()
			if err != nil {
				return nil, err
			}
			if n2 > n1 {
				stack.push(1)
			} else {
				stack.push(0)
			}

		case 'A':
			n1, n2, err := stack.popN2()
			if err != nil {
				return nil, err
			}
			if n1 != 0 && n2 != 0 {
				stack.push(1)
			} else {
				stack.push(0)
			}

		case 'O':
			n1, n2, err := stack.popN2()
			if err != nil {
				return nil, err
			}
			if n1 != 0 || n2 != 0 {
				stack.push(1)
			} else {
				stack.push(0)
			}

		case '!':
			n, err := stack.popN()
			if err != nil {
				return nil, err
			}

			if n == 0 {
				stack.push(1)
			} else {
				stack.push(0)
			}

		case '~':
			n, err := stack.popN()
			if err != nil {
				return nil, err
			}
			stack.push(^n)

		case 'i':
			if len(params) < 2 {
				return nil, errors.New("&i")
			}

			n1, ok := params[0].(int)
			if !ok {
				return nil, errors.New("&i")
			}

			n2, ok := params[1].(int)
			if !ok {
				return nil, errors.New("&i")
			}

			params[0] = n1 + 1
			params[1] = n2 + 1

		case '?':
			// nop
		case 't':
			n, err := stack.popN()
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
			format, c, err := parseFormat(b, in)
			if err != nil {
				return nil, err
			}

			if c == 's' {
				s, err := stack.popS()
				if err != nil {
					return nil, err
				}
				s = fmt.Sprintf(format, s)
				out.WriteString(s)
			} else {
				n, err := stack.popN()
				if err != nil {
					return nil, err
				}
				s := fmt.Sprintf(format, n)
				out.WriteString(s)
			}
		}
	}
}

func parseFormat(b byte, in io.ByteReader) (string, byte, error) {
	buf := new(bytes.Buffer)
	buf.WriteByte('%')

	var err error
	for {
		switch b {
		case ':':
			if buf.Len() > 0 {
				return "", 0, errors.New("illegal format")
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
