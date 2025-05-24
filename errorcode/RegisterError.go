package errorcode

import (
	"fmt"
	"log"
	"os"
)

// RegisterError same as before, adds to registry
func RegisterError(description string) {
	code, pkg, method, linefile := generateCodeAndInfo()

	Mu.Lock()
	Registry[code] = &ErrorEntry{
		Code:        code,
		Package:     pkg,
		Method:      method,
		Description: description,
		LineFile:    linefile,
	}

	ErrChan <- Registry

	Counters[codeKey(pkg, method)]++ // Increment count per block, adjust if needed
	Mu.Unlock()

	fmt.Printf("[%s] %s", code, description)
}

func WriteintoString() {
	for _, v := range <-ErrChan {
		lErrString := fmt.Sprintf("%s @@ %s @@ %s @@ %s @@ %s \n", v.Code, v.Package, v.Method, v.LineFile, v.Description)
		lfile, err := os.OpenFile("./log/elog.txt", os.O_RDWR|os.O_APPEND, 0666)
		log.Println("err : ", err)
		lfile.Write([]byte(lErrString))
	}
}
