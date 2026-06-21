/*
The MIT License (MIT)

Copyright © 2025 John, Sing Dao, Siu <john.sd.siu@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package lib

import (
	"slices"
	"strconv"
	"strings"

	"github.com/J-Siu/go-helper/v2/ezlog"
)

const VerDelimiters = "._-"

// return v1 > v2
func VerNewer(v1, v2 string) (newer bool) { return segmentCompare(segmentSplit(v1), segmentSplit(v2)) }

// return s1 > s2
func segmentCompare(s1, s2 []string) (newer bool) {
	prefix := "SegmentCompare"
	var (
		e     error
		equal = true
		int1  int
		int2  int
		len1  = len(s1)
		len2  = len(s2)
	)
	for i := range min(len1, len2) {
		ezlog.Debug().N(prefix).N("s1").N(i).M(s1[i]).Out()
		ezlog.Debug().N(prefix).N("s2").N(i).M(s2[i]).Out()
		int1, e = strconv.Atoi(trimNonNumeric(s1[i]))
		if e == nil {
			int2, e = strconv.Atoi(trimNonNumeric(s2[i]))
		}
		if e == nil {
			equal = equal && int1 == int2
			newer = int1 > int2
		}
		if e != nil || !equal {
			break
		}
	}

	if e == nil && equal {
		newer = len1 > len2
	}
	return newer
}

// Trim leading 'v', split s by VerDelimiters
func segmentSplit(s string) (segments []string) {
	var (
		char = ""
	)
	s = strings.TrimPrefix(s, "v")
	s = strings.TrimPrefix(s, "V")
	segments = []string{}
	for l := range len(s) {
		char = string(s[l])
		if strings.Contains(VerDelimiters, char) {
			segments = append(segments, "")
		} else {
			if len(segments) == 0 {
				segments = append(segments, "")
			}
			segments[len(segments)-1] += char
		}
	}
	return segments
}

func trimNonNumeric(s string) string {
	prefix := "trimNonNumeric"
	var (
		start int
		end   int
	)
	// get start
	for i, c := range s {
		start = i
		if c >= '0' && c <= '9' {
			break
		}
	}
	// get end
	for i, c := range slices.Backward([]byte(s)) {
		end = i
		if c >= '0' && c <= '9' {
			break
		}
	}
	end++

	ezlog.Debug().N(prefix).N("s").M(s).N("start").M(start).N("end").M(end).Out()
	ezlog.Debug().N(prefix).N("out").M(s[start:end]).Out()
	return s[start:end]
}
