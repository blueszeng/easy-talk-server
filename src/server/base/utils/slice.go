package utils

type SliceHandler func(from int, to int)

func Slice(total int, per int, h SliceHandler) {
	if total <= 0 {
		InvalidValueErr("Slice", "total <= 0")
		return
	}

	if per <= 0 {
		InvalidValueErr("Slice", "per <= 0")
		return
	}

	if h == nil {
		InvalidValueErr("Slice", "h == nil")
		return
	}

	from := 0
	to := 0
	times := (total-1)/per + 1
	for i := 0; i < times; i++ {
		from = to
		to = (i + 1) * per
		if to > total {
			to = total
		}
		h(from, to)
	}
}
