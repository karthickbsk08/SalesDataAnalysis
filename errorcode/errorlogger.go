package errorcode

/*
type ErrorEntry struct {
	Code        string
	Package     string
	Method      string
	Description string
	LineFile    string
}

var (
	mu       sync.Mutex
	registry = map[string]*ErrorEntry{}
	counters = map[string]int{}
	capRegex = regexp.MustCompile(`[A-Z]`)
)

// GenerateErrorCode and gather info from caller stack
func generateCodeAndInfo() (code, pkg, method, linefile string) {
	// Caller(2) = the function that called RegisterError
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return "UNKNOWN000", "unknown", "unknown", ""
	}
	fn := runtime.FuncForPC(pc)
	fullName := fn.Name() // e.g. github.com/user/project/payment.MakePayment

	pkg = extractPackageName(fullName)
	method = extractMethodName(fullName)

	pkgFirstChar := "X"
	if len(pkg) > 0 {
		pkgFirstChar = strings.ToUpper(pkg[:1])
	}

	caps := strings.Join(capRegex.FindAllString(method, -1), "")
	if caps == "" {
		caps = strings.ToUpper(method)
	}

	key := pkgFirstChar + caps

	mu.Lock()
	counters[key]++
	count := counters[key]
	mu.Unlock()

	code = fmt.Sprintf("%s%03d", key, count)
	linefile = fmt.Sprintf("%s:%d", file[strings.LastIndex(file, "/")+1:], line)
	return
}

// RegisterError registers an error and returns an error with code in the message
func RegisterError(description string) error {
	code, pkg, method, linefile := generateCodeAndInfo()

	mu.Lock()
	registry[code] = &ErrorEntry{
		Code:        code,
		Package:     pkg,
		Method:      method,
		Description: description,
		LineFile:    linefile,
	}
	mu.Unlock()

	return fmt.Errorf("[%s] %s", code, description)
}

// ExportCSV exports all registered errors to CSV with headers
func ExportCSV(path string) error {
	mu.Lock()
	defer mu.Unlock()

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	headers := []string{"Error Code", "Package", "Method", "Description", "Line/File"}
	if err := w.Write(headers); err != nil {
		return err
	}

	for _, e := range registry {
		record := []string{e.Code, e.Package, e.Method, e.Description, e.LineFile}
		if err := w.Write(record); err != nil {
			return err
		}
	}
	return nil
}

// Helpers to extract package and method from full func name
func extractPackageName(fullName string) string {
	// Try to get the last part of the package path before the method
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
*/
