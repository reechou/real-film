package ext

import (
	"testing"
	"fmt"
)

func TestFilmHttpClient(t *testing.T) {
	f := NewFilmHttpClient()
	fp, err := f.GetPlayer("acku", "CODUzODAxMg==~8c2ad552.acku")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(fp)
}
