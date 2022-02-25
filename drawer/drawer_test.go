package drawer

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestDrawText(t *testing.T) {
	images := []string{
		"image-320.jpg",
		"image-640.jpg",
		"image-1024.jpg",
		"image-2048.jpg",
	}

	for _, image := range images {
		buf, err := os.ReadFile(fmt.Sprintf("./fixtures/%s", image))
		if err != nil {
			t.Fatalf("error reading fixture: %v", err)
		}

		cases := map[string]string{
			"phrase-small":  "три",
			"phrase-mid-1":  "десять букв",
			"phrase-mid-2":  "средняя по протяженности",
			"phrase-mid-3":  "средняя по протяженности фраза",
			"phrase-long-1": "длинная по протяженности фраза с добавкой",
			"phrase-long-2": "длинная по протяженности фраза с добавкой и еще немного",
		}

		for k, v := range cases {
			out, err := DrawText(bytes.NewReader(buf), "../public/Lobster-Regular.ttf", v)
			if err != nil {
				t.Fatalf("error drawing text: %v", err)
			}

			output := new(bytes.Buffer)
			output.ReadFrom(out)

			folder := strings.Replace(image, ".jpg", "", -1)
			filename := fmt.Sprintf("./fixtures/%s-%s.png", folder, k)

			if err := os.WriteFile(filename, output.Bytes(), 0666); err != nil {
				t.Fatalf("error writing file: %v", err)
			}
		}
	}

}
