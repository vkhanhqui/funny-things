package utils

import "path/filepath"

func IndexPath(partition string) string {
	// Assuming the index files are stored in a directory named after the partition
	return filepath.Join(partition)
}
