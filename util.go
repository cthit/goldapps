package goldapps

import (
	"bytes"
	"fmt"
)

func printProgress(done int, total int) {
	p := (done * 100) / total
	builder := bytes.Buffer{}
	for i := 0; i < 100; i++ {
		if i < p {
			builder.WriteByte('=')
		} else if i == p {
			builder.WriteByte('>')
		} else {
			builder.WriteByte(' ')
		}

	}
	fmt.Printf("\rProgress: [%s] %d/%d", builder.String(), done, total)
	if done == total {
		fmt.Printf("\rDone\n")
	}
}
