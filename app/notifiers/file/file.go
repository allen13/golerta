package file

import (
	"fmt"
	"github.com/allen13/golerta/app/models"
	"io"
	"os"
)

type File struct {
	Files        []string `mapstructure:"file"`
	EnabledField bool     `mapstructure:"enabled"`
	writer       io.Writer
	closers      []io.Closer
}

func (f *File) Enabled() bool {
	return f.EnabledField
}

func (f *File) Trigger(alert models.Alert) error {
	alertMessage := "triggering alert: " + alert.String()
	return f.writeMessage(alertMessage)
}

func (f *File) Acknowledge(alert models.Alert) error {
	alertMessage := "acknowledging alert: " + alert.String()
	return f.writeMessage(alertMessage)
}

func (f *File) Resolve(alert models.Alert) error {
	alertMessage := "resolving alert: " + alert.String()
	return f.writeMessage(alertMessage)
}

func (f *File) writeMessage(message string) error {
	_, err := f.writer.Write([]byte(message + "\n"))
	if err != nil {
		return fmt.Errorf("FAILED to write message: %s, %s", message, err)
	}

	return nil
}

func (f *File) Init() error {
	writers := []io.Writer{}

	if len(f.Files) == 0 {
		f.Files = []string{"stdout"}
	}

	for _, file := range f.Files {
		if file == "stdout" {
			writers = append(writers, os.Stdout)
			f.closers = append(f.closers, os.Stdout)
		} else {
			var of *os.File
			var err error
			if _, err := os.Stat(file); os.IsNotExist(err) {
				of, err = os.Create(file)
			} else {
				of, err = os.OpenFile(file, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
			}

			if err != nil {
				return err
			}
			writers = append(writers, of)
			f.closers = append(f.closers, of)
		}
	}
	f.writer = io.MultiWriter(writers...)
	return nil
}

func (f *File) Close() error {
	var errS string
	for _, c := range f.closers {
		if err := c.Close(); err != nil {
			errS += err.Error() + "\n"
		}
	}
	if errS != "" {
		return fmt.Errorf(errS)
	}
	return nil
}
