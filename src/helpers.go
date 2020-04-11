package cache


func countValues(a []int, i int) int {
    count := 0
    for _, val := range a {
        if val == i {
            count++
        }
    }
    return count
}

// JOIN
func joinInts(a []int, b []int) []int {
    var res []int
    for _, val := range a {
        res = append(res, val)
    }
    for _, val := range b {
        res = append(res, val)
    }
    return res
}


// Equal tells whether a and b contain the same elements.
// A nil argument is equivalent to an empty slice.
func IntArrayEqual(a, b []int) bool {
    if len(a) != len(b) {
        return false
    }
    for i, v := range a {
        if v != b[i] {
            return false
        }
    }
    return true
}


func max(a int, b int) int {
    if a > b {
        return a
    }
    return b
}

func min(a int, b int) int {
    if a < b {
        return a
    }
    return b
}



