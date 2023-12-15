package store

type Store interface {
	SetValue(newVal string, key string)
	GetValue(key string) (string, int)
	NotifyValue(curVal string, key string, curGeneration int) bool
}
