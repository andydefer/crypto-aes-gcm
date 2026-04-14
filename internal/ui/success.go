package ui

import (
	"fmt"

	"github.com/andydefer/crypto-aes-gcm/internal/utils"
)

// PrintSuccess displays completion information with file size.
func PrintSuccess(output string, size int64) {
	fmt.Println()
	SuccessColor.Printf("✅ Operation successful!\n")
	InfoColor.Printf("📄 Output: %s\n", output)
	InfoColor.Printf("📏 Size:   %s\n", utils.FormatFileSize(size))
	fmt.Println()
}
