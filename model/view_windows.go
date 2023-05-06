package msi

var (
	procMsiViewExecute = modMsi.NewProc("MsiViewExecute")
)

type View struct {
	handle uintptr
}

func (v *View) Close() error {
	if v.handle == 0 {
		return nil
	}

	if ret, _, err := procMsiCloseHandle.Call(v.handle); ret != 0 {
		return err
	}

	v.handle = 0
	return nil
}

func (v View) Execute(r Record) (int, error) {
	ret, _, err := procMsiViewExecute.Call(
		v.handle,
		r.handle,
	)
	return int(ret), err
}
