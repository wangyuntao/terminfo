# terminfo

A terminfo library for Go

## Example

```go
package main

import (
        "fmt"
        "log"

        "github.com/wangyuntao/terminfo"
)

func main() {
        ti, err := terminfo.LoadEnv()
        failIfErr(err)

        err = ti.ClearScreen()
        failIfErr(err)

        rainbow(ti, "HELLO", terminfo.ColorRed)
        err = ti.CursorAddress(1, 5)
        failIfErr(err)

        rainbow(ti, "WORLD", terminfo.ColorRed)
        fmt.Println()
}

func rainbow(ti *terminfo.Terminfo, s string, c int) {
        for _, r := range s {
                err := ti.ColorFg(c)
                failIfErr(err)

                fmt.Printf("%c", r)
                c++
        }
}

func failIfErr(err error) {
        if err != nil {
                log.Fatal(err)
        }
}
```


## References

- [ANSI escape code][l1]
- [term(5)][l2]
- [terminfo(5)][l3]
- [ncurses source code][l4]

<!-- references -->

[l1]: https://en.wikipedia.org/wiki/ANSI_escape_code
[l2]: https://invisible-island.net/ncurses/man/term.5.html
[l3]: https://invisible-island.net/ncurses/man/terminfo.5.html
[l4]: https://invisible-island.net/datafiles/current/ncurses.tar.gz