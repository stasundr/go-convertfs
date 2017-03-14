package admutils

// HashIt takes a string as argument
// and calculates hash for it
// (see https://github.com/DReichLab/EIG/blob/master/src/admutils.c#L684)
func HashIt(str string) uint32 {
	var hash, length uint32 = 0, uint32(len(str))

	for j := uint32(0); j < length; j++ {
		hash *= 23
		hash += uint32(str[j])
	}

	return hash
}

// HashArr takes an array of strings as argument
// and calculates hash for it
// (see https://github.com/DReichLab/EIG/blob/master/src/admutils.c#L667)
func HashArr(xarr []string) uint32 {
	var hash, thash, nxarr uint32
	nxarr = uint32(len(xarr))

	for i := uint32(0); i < nxarr; i++ {
		thash = HashIt(xarr[i])
		hash *= 17
		hash ^= thash
	}

	return hash
}
