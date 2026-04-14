// Package service provides business logic for encryption and decryption.
package service

import (
	"io"
	"os"

	"github.com/andydefer/crypto-aes-gcm/internal/ui"
	"github.com/andydefer/crypto-aes-gcm/pkg/cryptolib"
)

// ExecuteEncryption performs the encryption operation.
func ExecuteEncryption(input, output, password string, workerCount int, quiet bool) error {
	fileInfo, err := os.Stat(input)
	if err != nil {
		return err
	}
	fileSize := fileInfo.Size()

	var bar ui.ProgressBar
	if !quiet {
		bar = ui.CreateProgressBar(fileSize, "🔒 Encrypting")
	} else {
		bar = &noopProgressBar{}
	}

	encryptor, err := cryptolib.NewEncryptor(workerCount)
	if err != nil {
		return err
	}

	reader := &progressReader{
		r:     mustOpenFile(input),
		bar:   bar,
		total: fileSize,
	}
	defer reader.Close()

	outFile, err := os.Create(output)
	if err != nil {
		return err
	}
	defer outFile.Close()

	if err := encryptor.Encrypt(reader, outFile, password); err != nil {
		_ = bar.Clear()
		return err
	}

	_ = bar.Finish()
	ui.PrintSuccess(output, fileSize)
	return nil
}

// progressReader wraps an io.ReadCloser to track reading progress.
type progressReader struct {
	r     io.ReadCloser
	bar   ui.ProgressBar
	total int64
	read  int64
}

func (pr *progressReader) Read(p []byte) (n int, err error) {
	n, err = pr.r.Read(p)
	pr.read += int64(n)
	_ = pr.bar.Set64(pr.read)
	return n, err
}

func (pr *progressReader) Close() error {
	return pr.r.Close()
}

type noopProgressBar struct{}

func (n *noopProgressBar) Set64(int64) error { return nil }
func (n *noopProgressBar) Finish() error     { return nil }
func (n *noopProgressBar) Clear() error      { return nil }

func mustOpenFile(path string) io.ReadCloser {
	f, err := os.Open(path)
	if err != nil {
		ui.ErrorColor.Fprintf(os.Stderr, "❌ Erreur: %v\n", err)
		os.Exit(1)
	}
	return f
}
