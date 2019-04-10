package terminfo

func (ti *Terminfo) Cursor(row, col int) error {
	return ti.Do(CursorAddress, row, col)
}
