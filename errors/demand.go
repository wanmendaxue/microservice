package errors

import (
	"fmt"
)

type demand string
func (d demand) Error() string {
	return string(d)
}

func NewDemandError(msg string) error {
	return demand(msg)
}

func NewDemandErrorf(format string, a ...interface{}) error {
	return NewDemandError(fmt.Sprintf(format, a...))
}

func (d demand) Demand() string {
	return string(d)
}