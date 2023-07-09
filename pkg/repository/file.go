package repository

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
)

type File struct {
	baseDir string
}

func NewFile(baseDir string) *File {
	return &File{
		baseDir: baseDir,
	}
}

func (f *File) calculatePathForKey(key string) string {
	return path.Join(f.baseDir, strings.ReplaceAll(key, "/", "?"))
}

func (f *File) StoreValue1000(ctx context.Context, key string, value1000 int) error {
	filename := f.calculatePathForKey(key)

	var data bytes.Buffer

	fmt.Fprintf(&data, "%d", value1000)

	err := os.WriteFile(filename, data.Bytes(), 0644)

	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (f *File) GetValue1000(ctx context.Context, key string) (int, error) {
	filename := f.calculatePathForKey(key)

	contents, err := os.ReadFile(filename)

	if err != nil {
		if os.IsNotExist(err) {
			return 0, fmt.Errorf("not found")
		}

		return 0, fmt.Errorf("failed to read file: %w", err)
	}

	val, err := strconv.Atoi(string(contents))

	if err != nil {
		return 0, fmt.Errorf("failed to parse value %q: %w", string(contents), err)
	}

	return val, nil
}
