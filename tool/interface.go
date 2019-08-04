package tool

// ILogger provides logging capabilities.
type ILogger interface {
	Fatalf(format string, v ...interface{})
	Printf(format string, v ...interface{})
}

// IStorage provides storage for sensitive data.
type IStorage interface {
	//Filename(params ...string) string
	LoadFile(v interface{}, params ...string) error
	Save(v interface{}, params ...string) error
}
