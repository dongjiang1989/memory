package Utils

import (
	"os"
	"path/filepath"
)

func Path() string {
	return os.Args[0]
}

// real path
func RealPath() string {
	path := os.Args[0]
	s, _ := filepath.EvalSymlinks(path)
	s, _ = filepath.Abs(s)
	return s
}

// abs path
func AbsPath() string {
	path := os.Args[0]
	s, _ := filepath.Abs(path)
	return s
}

//norm path
func NormPath() string {
	path := os.Args[0]
	return filepath.Clean(path)
}

//pwd path
func Getpwd() string {
	path := os.Args[0]
	s, _ := os.Getwd()
	return filepath.Join(s, path)
}

//dir name
func Dirname() string {
	dir, _ := filepath.Split(AbsPath())
	return dir
}

// bin file name
func Filename() string {
	_, filename := filepath.Split(AbsPath())
	return filename
}
