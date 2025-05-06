package internal

import "sync"

// See: https://github.com/valyala/fasthttp/issues/920
var state = sync.Map{}

func Set(key string, value any) {
	state.Store(key, value)
}

func MustGet(key string) any {
	if dep, ok := state.Load(key); ok {
		return dep
	}

	panic("state: dependency not found!")
}

func MustGetState[T any](key string) T {
	dep, ok := MustGet(key).(T)
	if !ok {
		panic("state: dependency not found!")
	}

	return dep
}
