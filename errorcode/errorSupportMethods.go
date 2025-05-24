package errorcode

import (
	"fmt"
	"runtime"
	"strings"
)

// generateCodeAndInfo returns code, package, method, linefile info
func generateCodeAndInfo() (string, string, string, string) {
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return "UNKNOWN000", "unknown", "unknown", ""
	}
	fn := runtime.FuncForPC(pc)
	fullName := fn.Name()

	pkg := extractPackageName(fullName)
	method := extractMethodName(fullName)

	pkgFirstChar := "X"
	if len(pkg) > 0 {
		pkgFirstChar = strings.ToUpper(pkg[:1])
	}
	caps := extractCaps(method)
	key := pkgFirstChar + caps

	Mu.Lock()
	count := Counters[key] + 1
	Counters[key] = count
	Mu.Unlock()

	code := fmt.Sprintf("%s%03d", key, count)
	linefile := fmt.Sprintf("%s:%d", file[strings.LastIndex(file, "/")+1:], line)
	return code, pkg, method, linefile
}

func codeKey(pkg, method string) string {
	caps := extractCaps(method)
	return strings.ToUpper(pkg[:1]) + caps
}

func extractCaps(s string) string {
	var caps []rune
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			caps = append(caps, r)
		}
	}
	if len(caps) == 0 {
		return strings.ToUpper(s)
	}
	return string(caps)
}

func extractPackageName(fullName string) string {
	lastSlash := strings.LastIndex(fullName, "/")
	if lastSlash == -1 {
		lastSlash = 0
	} else {
		lastSlash++
	}
	dot := strings.Index(fullName[lastSlash:], ".")
	if dot == -1 {
		return ""
	}
	return fullName[lastSlash : lastSlash+dot]
}

func extractMethodName(fullName string) string {
	lastDot := strings.LastIndex(fullName, ".")
	if lastDot == -1 || lastDot == len(fullName)-1 {
		return fullName
	}
	return fullName[lastDot+1:]
}
