package countingwriter

import (
	"fmt"
	"os"
	"text/template"
)

func ExampleCountingWriter_BytesWritten() {

	countingWriter := New(os.Stdout)
	t := template.Must(template.New("t").Parse("Hello {{.Name}}\n"))
	_ = t.Execute(countingWriter, struct{ Name string }{Name: "JF"})
	bytes := countingWriter.BytesWritten()
	fmt.Printf("Bytes written: %v", bytes)

	// Output:
	// Hello JF
	// Bytes written: 9
}
