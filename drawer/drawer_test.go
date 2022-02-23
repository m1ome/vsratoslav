package drawer

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

func TestDrawText(t *testing.T) {
	buf, err := os.ReadFile("./fixtures/image.jpg")
	if err != nil {
		t.Fatalf("error reading fixture: %v", err)
	}

	cases := map[string]string{
		"small":  "лох",
		"middle": "Ща будут пиздосики",
		"long":   "Типичная русская семья времен кровавого царизма",
	}

	for k, v := range cases {
		out, err := DrawText(bytes.NewReader(buf), "../Lobster-Regular.ttf", v)
		if err != nil {
			t.Fatalf("error drawing text: %v", err)
		}

		output := new(bytes.Buffer)
		output.ReadFrom(out)

		if err := os.WriteFile(fmt.Sprintf("./fixtures/%s.png", k), output.Bytes(), 0666); err != nil {
			t.Fatalf("error writing file: %v", err)
		}
	}

}
