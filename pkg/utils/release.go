package utils

type Releasable interface {
	Release() error
}

func Release(releasables ...Releasable) error {
	var err error
	for _, releasable := range releasables {
		if er := releasable.Release(); er != nil {
			if err == nil {
				err = er
			}
		}
	}
	return err
}
