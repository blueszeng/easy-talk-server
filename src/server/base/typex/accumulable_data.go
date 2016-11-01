package typex

////////////////////////////////////////////
// type, const, var
//
type AccumulableBoolData struct {
	data int
}
type AccumulableMaxMinData map[int32]int

////////////////////////////////////////////
// func
//

//
// AccumulableBoolData method
//
func (d *AccumulableBoolData) Is() bool {
	return d.data > 0
}

func (d *AccumulableBoolData) Set(is bool) {
	if is {
		d.data += 1
	} else {
		d.data -= 1
		if d.data < 0 {
			d.data = 0
		}
	}
}

//
// AccumulableMaxMinData method
//
func (d *AccumulableMaxMinData) Max() int32 {
	if len(*d) == 0 {
		return 0
	}

	ret := int32(0)
	for v, num := range *d {
		if num > 0 {
			if ret == 0 {
				ret = v
			} else if v > ret {
				ret = v
			}
		}
	}
	return ret
}

func (d *AccumulableMaxMinData) Min() int32 {
	if len(*d) == 0 {
		return 0
	}

	ret := int32(0)
	for v, num := range *d {
		if num > 0 {
			if ret == 0 {
				ret = v
			} else if v < ret {
				ret = v
			}
		}
	}
	return ret
}

func (d *AccumulableMaxMinData) Add(v int32) {
	if v <= 0 {
		return
	}

	(*d)[v] += 1
}

func (d *AccumulableMaxMinData) Rmv(v int32) {
	if len(*d) == 0 {
		return
	}

	if (*d)[v] > 1 {
		(*d)[v] -= 1
	} else {
		delete(*d, v)
	}
}

func (d *AccumulableMaxMinData) Len() int {
	return len(*d)
}
