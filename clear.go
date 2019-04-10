package terminfo

func (ti *Terminfo) ClearScreen() error {
	return ti.Do(ClearScreen)
}
