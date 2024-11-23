package bloomfilter

import "errors"

type BitList interface {
	Set(index uint64) error
	Check(index uint64) (bool, error)
	Len() uint64
	SetCapacity(bits uint64)
}

type Int64Slice []int64

func (i *Int64Slice) SetCapacity(bits uint64) {
	*i = make([]int64, (bits+63)/64)
}

func (i *Int64Slice) Set(index uint64) error {
	if index >= i.Len() {
		return errors.New("wrong index")
	}
	(*i)[index] |= 1 << (index % 64)
	return nil
}

func (i *Int64Slice) Check(index uint64) (bool, error) {
	if index >= i.Len() {
		return false, errors.New("wrong index")
	}
	return !((*i)[index]&(1<<(index%64)) == 0), nil
}

func (i *Int64Slice) Len() uint64 {
	return uint64(len(*i))
}
