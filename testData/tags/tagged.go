//go:build tagged
// +build tagged

package tags

// Tagged is only visible with tags.
func Tagged() int {
	return 7
}
