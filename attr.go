package terminfo

func (ti *Terminfo) EnterStandoutMode() error {
	return ti.Do(EnterStandoutMode)
}

func (ti *Terminfo) EnterUnderlineMode() error {
	return ti.Do(EnterUnderlineMode)
}

func (ti *Terminfo) EnterReverseMode() error {
	return ti.Do(EnterReverseMode)
}

func (ti *Terminfo) EnterBlinkMode() error {
	return ti.Do(EnterBlinkMode)
}

func (ti *Terminfo) EnterDimMode() error {
	return ti.Do(EnterDimMode)
}

func (ti *Terminfo) EnterBoldMode() error {
	return ti.Do(EnterBoldMode)
}

func (ti *Terminfo) EnterItalicsMode() error {
	return ti.Do(EnterItalicsMode)
}

func (ti *Terminfo) ExitAttributeMode() error {
	return ti.Do(ExitAttributeMode)
}
