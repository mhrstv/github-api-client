package input

import (
	"bufio"
	"os"
	"strings"
)

func ReadUsernames(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var users []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		if name := strings.TrimSpace(sc.Text()); name != "" {
			users = append(users, name)
		}
	}
	return users, sc.Err()
}
