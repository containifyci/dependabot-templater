package search

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUniqueSlice(t *testing.T) {
	t.Parallel()

	foundFiles := UniqueStringSlice{
		unqiue: make(map[string]bool),
	}
	foundFiles.Add("test")
	foundFiles.Add("../../../test")
	assert.Equal(t, 1, len(foundFiles.unqiue))
}

// func TestNormalizePath(t *testing.T) {
// 	t.Parallel()

// 	assert.Equal(t, ".", NormalizePath("test"))
// 	assert.Equal(t, ".", NormalizePath("./test"))
// 	assert.Equal(t, "test", NormalizePath("../project/subdir/test"))
// 	assert.Equal(t, "test", NormalizePath("../../project/subdir/test"))
// }

func TestSearchForFiles(t *testing.T) {
	t.Parallel()

	foundFiles, err := SearchForFiles("./test_path", "requirements.txt")
	assert.NoError(t, err)
	assert.Equal(t, []string{"test_path/projecte"}, foundFiles)
}

func TestSearchForString(t *testing.T) {
	t.Parallel()

	foundFiles, err := SearchForString("./test_path", "terraform")
	assert.NoError(t, err)
	assert.Equal(t, []string{"test_path/projectb"}, foundFiles)
}

func TestSearchForFolder(t *testing.T) {
	t.Parallel()

	foundFiles, err := SearchForFolder("./test_path", "projectc")
	assert.NoError(t, err)
	assert.Equal(t, []string{"test_path/projectc"}, foundFiles)
}

func TestSearchNoResult(t *testing.T) {
	t.Parallel()

	foundFiles, err := SearchForFiles("./test_path", "nofile")
	assert.NoError(t, err)
	assert.Equal(t, nilSlice, foundFiles)

	foundFiles, err = SearchForFolder("./test_path", "nofolder")
	assert.NoError(t, err)
	assert.Equal(t, nilSlice, foundFiles)

	foundFiles, err = SearchForString("./test_path", "nomatch")
	assert.NoError(t, err)
	assert.Equal(t, nilSlice, foundFiles)
}

// test utility

var nilSlice []string
