package utils

import "errors"
func Map[I any, O any](in []I, f func(i I) (o O, err error)) ([]O, error) {
	out := make([]O, len(in))
	var ERR error
	k := 0
	for _, i := range in {
		o, err := f(i)
		if err != nil {
			ERR = errors.Join(ERR, err)
			continue
		}
		out[k] = o
		k++
	}
	return out[:k], ERR
}
func MapIdx[I any, O any](in []I, f func(i I, idx int) (o O, err error)) ([]O, error) {
	out := make([]O, len(in))
	var ERR error
	k := 0
	for idx, i := range in {
		o, err := f(i, idx)
		if err != nil {
			ERR = errors.Join(ERR, err)
			continue
		}
		out[k] = o
		k++
	}
	return out[:k], ERR
}
