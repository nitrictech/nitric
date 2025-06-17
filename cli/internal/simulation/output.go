package simulation

import (
	"fmt"
	"io"
	"strings"
)

type PrefixWriter struct {
	writer io.Writer
	prefix string
}

func (p *PrefixWriter) Write(content []byte) (int, error) {
	value := strings.TrimSuffix(string(content), "\n")

	split := strings.Split(value, "\n")
	value = strings.Join(split, "\n"+p.prefix) + "\n"

	_, err := fmt.Fprintf(p.writer, "%s%s", p.prefix, value)
	if err != nil {
		return 0, err
	}

	return len(content), nil
}

func NewPrefixWriter(prefix string, writer io.Writer) *PrefixWriter {
	return &PrefixWriter{
		prefix: prefix,
		writer: writer,
	}
}
