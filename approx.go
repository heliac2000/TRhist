package main

type ApproxRepeat struct {
	Repeat string
	Length int
}

type RepeatAttribute struct {
	Start, End, Unit, Length int
	Repeat                   string
	RepeatLength             int
}

func ExtendRepeat(read string, rm map[Result]bool) []ApproxRepeat {
	max, rr := 0, make([]Result, 0)
	for v := range rm {
		if v.Length < 3 {
			continue
		}
		if l := v.Length * v.Count; l >= max {
			if l > max {
				max, rr = l, nil
			}
			rr = append(rr, v)
		}
	}

	if max < 8 {
		return nil
	}

	apps, rl := make([]ApproxRepeat, 0), len(read)
	for _, r := range rr {
		rep := read[r.From:(r.From + r.Length)]
		right := ExtendRight(read, rep, r.From+(r.Length*r.Count), 0, r.Length, rl)
		left := ExtendLeft(read, rep, 0, r.From, r.Length, rl)
		if left < r.From || right > r.From+r.Length*r.Count {
			delete(rm, r)
			apps = append(apps, ApproxRepeat{
				Repeat: read[r.From:(r.From + r.Length)],
				Length: right - left + 1,
			})
		}
	}

	return apps
}

func ExtendRight(read, rep string, start, end, unit, rl int) int {
	for start < rl {
		end = start + unit
		if d := end - rl; d > 1 {
			break
		} else if d > 0 {
			end = rl
		}
		if rep == read[start:end] ||
			extend1Mismatch(read, rep, start, end) {
			start = end
		} else if (end+1) <= rl &&
			extend1Insertion(read, rep, unit, start, end+1) {
			start = end + 1
		} else if extend1Deletion(read, rep, start, end-1) {
			start = end - 1
		} else {
			break
		}
	}

	return (start - 1)
}

func ExtendLeft(read, rep string, start, end, unit, rl int) int {
	for end > 0 {
		if start = end - unit; start < -1 {
			break
		} else if start < 0 {
			start = 0
		}
		if rep == read[start:end] ||
			extend1Mismatch(read, rep, start, end) {
			end = start
		} else if start > 0 &&
			extend1Insertion(read, rep, unit, start-1, end) {
			end = start - 1
		} else if extend1Deletion(read, rep, start+1, end) {
			end = start + 1
		} else {
			break
		}
	}

	return end
}

func extend1Mismatch(read, rep string, start, end int) bool {
	for i, j, mis := start, 0, 0; i < end; i, j = i+1, j+1 {
		if read[i] != rep[j] {
			if mis++; mis > 1 {
				return false
			}
		}
	}
	return true
}

func extend1Insertion(read, rep string, unit, start, end int) bool {
	for i, j, ins := start, 0, 0; i < end && j < unit; i, j = i+1, j+1 {
		if read[i] != rep[j] {
			if ins++; ins > 1 {
				return false
			}
			j--
		}
	}
	return true
}

func extend1Deletion(read, rep string, start, end int) bool {
	for i, j, del := start, 0, 0; i < end; i, j = i+1, j+1 {
		if read[i] != rep[j] {
			if del++; del > 1 {
				return false
			}
			i--
		}
	}
	return true
}
