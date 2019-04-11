package terminfo

func (ti *Terminfo) CursorAddress(row, col int) error {
	return ti.Do(CursorAddress, row, col)
}
