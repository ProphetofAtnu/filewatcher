package filewatch

import (
	"crypto/md5"
	"errors"
	"io"
	"log"
	"os"
	"time"
)

// Constants defining watcher states
const (
	// Watcher is actively watching the file
	Active = iota
	// Watcher is active with errors
	Errors
	// Changed with errors
	ErrorsChanged
	// Watcher has detected a change on this interval
	Changed
	// Watcher has been manually halted
	Stopped
	// Watcher has been stopped by recoverable errors
	Fault
	// Watcher cannot continue
	Panic
)

type ErrorState struct {
	Error       error
	Recoverable bool
}

// Define errors for the watcher error stack
var (
	ErrOpen          = ErrorState{errors.New("Error opening file"), false}
	ErrFileLost      = ErrorState{errors.New("File Lost"), true}
	ErrFileNotExist  = ErrorState{errors.New("File no longer exists"), false}
	ErrCannotProcess = ErrorState{errors.New("Processing Failed"), true}
)

type ErrorStack []ErrorState

func (es *ErrorStack) Clear(logger *log.Logger) (ok bool) {
	var nes ErrorStack
	for _, e := range *es {
		if e.Recoverable {
			nes = append(nes, e)
		}
	}
}

// State defines the state of the Watcher
type State uint8

// Watcher defines a structure that holds the file to be watched
type Watcher struct {
	File     string
	State    State
	Interval time.Duration
	Modified time.Time
	Hash     []byte
	Changed  chan State
}

func (w *Watcher) Watch() {
	go func() {
		for {
		}
	}()
}

func (w Watcher) checkModified() error {
	s, err := os.Lstat(w.File)
	if err != nil {
		return err
	}
	return nil
}

func NewWatcher(filename string, interval time.Duration) (*Watcher, error) {
	stat, err := os.Lstat(filename)
	if err != nil {
		return nil, err
	}
	modtime := stat.ModTime()

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	hasher := md5.New()
	io.Copy(hasher, file)
	hash := hasher.Sum(nil)

	return &Watcher{
		filename,
		Stopped,
		interval,
		modtime,
		hash,
		make(chan State, 1),
	}, nil

}
