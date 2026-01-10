package database

type KeyRange []string

func (kr *KeyRange) Pop() (string, bool) {
	if len(*kr) == 0 {
		return "", false
	}
	first := (*kr)[0]
	*kr = (*kr)[1:]
	return first, true
}

func (kr KeyRange) Empty() bool {
	return len(kr) == 0
}

type KeyPool interface {
	ReserveRange() (KeyRange, error)
}
