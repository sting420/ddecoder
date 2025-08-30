package decoder

import (
	"bytes"
	"fmt"
	"strings"
)

type OrderedMap map[string]string

func (om OrderedMap) ToJson(order ...string) string {
	buf := &bytes.Buffer{}
	buf.Write([]byte{'{', '\n'})
	l := len(order)
	for i, k := range order {
		fmt.Fprintf(buf, "\t\"%s\": \"%v\"", strings.ReplaceAll(k, "\"", `\"`), strings.ReplaceAll(om[k], "\"", `\"`))
		if i < l-1 {
			buf.WriteByte(',')
		}
		buf.WriteByte('\n')
	}
	buf.Write([]byte{'}', '\n'})
	return buf.String()
}

func reverseCalcShift(val byte) int32 {
	integer := int32(val)
	var x int32 = 0
	if integer == 95 {
		x = 1
	} else if integer > 96 {
		x = integer - 59 // Corresponds to n > 37
	} else if integer > 64 {
		x = integer - 53 // Corresponds to n > 11
	} else if integer > 47 {
		x = integer - 46 // Corresponds to n > 1
	} else {
		x = (integer - 45) / 50 // Corresponds to n <= 1
	}
	return x
}

func splitChunk(r int32) []int32 {
	var originalBytes []int32

	originalBytes = append(originalBytes, (r>>16)&0xFF) // First byte
	originalBytes = append(originalBytes, (r>>8)&0xFF)  // Second byte
	originalBytes = append(originalBytes, (r & 0xFF))   // Third byte

	return originalBytes
}
func charCodeAt(s string, n int) int {
	i := 0
	for _, r := range s {
		if i == n {
			return int(r)
		}
		i++
	}
	return 0
}

func simpleHash(str string) int {
	if str == "" {
		return 1789537805
	}
	var t int = 0
	for e := 0; e < len(str); e++ {
		t = int(int32(t)<<5) - t + charCodeAt(str, e)
	}
	if t == 0 {
		return 1789537805
	} else {
		return t
	}

}
func secureAt(a []byte, i int) byte {
	if len(a) <= i {
		return 0
	} else {
		return a[i]
	}
}

func DecodePayload(payload string, cid string, hash string, scriptSeed int32) (string, error) {
	var seed int32 = int32(simpleHash(cid)) ^ int32(simpleHash(hash)) ^ scriptSeed

	byteArray := []byte(payload)
	firstStepArray := []int32{}
	for i := 0; i < len(byteArray); i += 4 {
		var chunk int32 = int32(reverseCalcShift(secureAt(byteArray, i))<<18) + int32(reverseCalcShift(secureAt(byteArray, i+1))<<12) + int32(reverseCalcShift(secureAt(byteArray, i+2))<<6) + int32(reverseCalcShift(secureAt(byteArray, i+3)))
		split := splitChunk(chunk)
		firstStepArray = append(firstStepArray, split...)
	}
	secondStepArray := []int32{}
	for i := 0; i < len(firstStepArray); i++ {

		var shift int32 = int32(16 - ((i % 3) * 8))
		secondStepArray = append(secondStepArray, ((seed>>shift)&255)^firstStepArray[i])
		if (i % 3) == 2 {
			seed ^= seed << 13
			seed ^= seed >> 17
			seed ^= seed << 5
		}

	}
	keys := []string{}
	vals := []string{}
	var i int
	for i = 0; i < len(secondStepArray); {
		if secondStepArray[i] == 123 || secondStepArray[i] == 44 {
			key := ""
			i++
			for ; i < len(secondStepArray); i++ {
				if secondStepArray[i] == 58 {

					break
				} else {
					key += string(rune(byte(secondStepArray[i])))
				}
			}
			keys = append(keys, key)
		}
		if secondStepArray[i] == 58 {
			val := ""
			i++
			inStr := false
			for ; i < len(secondStepArray); i++ {
				if secondStepArray[i] == 34 {
					inStr = !inStr
				}
				if secondStepArray[i] == 44 && !inStr {
					break
				} else {
					val += string(rune(secondStepArray[i]))
				}
			}
			vals = append(vals, val)
		}
	}

	str := "{"
	for i := 0; i < len(keys); i += 1 {
		if strings.Contains(keys[i], "r3n") {
			vals[i] = strings.Join(strings.SplitAfter(vals[i], `"`)[:len(strings.SplitAfter(vals[i], `"`))-1], "")
		}
		s := fmt.Sprintf("%s : %s,\n", keys[i], vals[i])
		str += s
	}
	return str + "}", nil
}
