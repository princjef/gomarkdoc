//go:build tagged
// +build tagged

package untagged

// Tagged is only visible with tags.
func Tagged() int {
	return 7
}
