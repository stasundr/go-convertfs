package admutils

// Hashit takes a string as argument
// and calculates hash for it
// (see https://github.com/DReichLab/EIG/blob/master/src/admutils.c :683)
func Hashit(str string) int {
	var hash, length = 0, len(str)

	for j := 0; j < length; j++ {
		hash *= 23
		hash += int(str[j])
	}

	return hash
}

// Hasharr takes an array of strings as argument
// and calculates hash for it
// (see https://github.com/DReichLab/EIG/blob/master/src/admutils.c :666)
func Hasharr(xarr []string) int {
	var hash, thash, nxarr int
	hash = 0
	nxarr = len(xarr)

	for i := 0; i < nxarr; i++ {
		thash = Hashit(xarr[i])
		hash *= 17
		hash ^= thash
	}

	return hash
}
