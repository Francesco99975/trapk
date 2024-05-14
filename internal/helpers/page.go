package helpers

import (
	"bytes"
	"context"

	"github.com/a-h/templ"
)

func GeneratePage(page templ.Component) ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	err := page.Render(context.Background(), buf)

	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}
