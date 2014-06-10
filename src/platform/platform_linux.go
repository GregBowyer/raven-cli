package platform

import (
	"syscall"
)

func charsToString(ca []int8) string {
	s := make([]byte, len(ca))
	var lens int
	for ; lens < len(ca); lens++ {
		if ca[lens] == 0 {
			break
		}
		s[lens] = uint8(ca[lens])
	}
	return string(s[0:lens])
}

func GetPlatform() Platform {
	var uname syscall.Utsname

	err := syscall.Uname(&uname)

	if err != nil {
		panic("Unable to provide a uname call for this platform")
	}

    return Platform {
        OSName: charsToString(uname.Sysname[:]),
        Release: charsToString(uname.Release[:]),
        Architecture: charsToString(uname.Machine[:]),
    }
}
