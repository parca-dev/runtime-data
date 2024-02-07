package libc

type Layout struct {
	pthreadSpecific1stblock uint64 `offsetof:"pthread.specific_1stblock" yaml:"pthread_specific_1stblock"`
	pthreadKeyData          uint64 `offsetof:"pthread_key_data.data"     yaml:"pthread_key_data"`
}

func (Layout) Data() ([]byte, error) {
	return nil, nil
}

type PThreadOffsets2 struct {
	pthread struct {
		specific_1stblock uint64 `offsetof:"specific_1stblock"`
	} `offsetof:"pthread"`

	pthread_key_data struct {
		data uint64 `offsetof:"data"`
	} `offsetof:"pthread_key_data"`
}
