package tool

// ILogger provides logging capabilities.
type ILogger interface {
	Fatalf(format string, v ...interface{})
	Printf(format string, v ...interface{})
}

// IStorage provides storage for sensitive data.
type IStorage interface {
	Load(v interface{}, params ...string) error
	Save(v interface{}, params ...string) error
	Delete(params ...string) error
	Filename(params ...string) string
}
