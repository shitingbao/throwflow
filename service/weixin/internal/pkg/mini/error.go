package mini

import (
	"fmt"
)

type MiniError struct {
	Description string
}

func (oe *MiniError) Error() string {
	return fmt.Sprintf(oe.Description)
}

func NewMiniError(description string) error {
	return &MiniError{
		Description: description,
	}
}

func (me *MiniError) GetDescription() string {
	return me.Description
}
