package service

import (
	"encoding/binary"
	"os"

	"github.com/andydefer/crypto-aes-gcm/internal/ui"
	"github.com/andydefer/crypto-aes-gcm/pkg/cryptolib"
)

// ExecuteDecryption performs the decryption operation.
func ExecuteDecryption(input, output, password string, quiet bool) error {
	fileInfo, err := os.Stat(input)
	if err != nil {
		return err
	}
	fileSize := fileInfo.Size()

	var bar ui.ProgressBar
	if !quiet {
		bar = ui.CreateProgressBar(fileSize, "🔓 Decrypting")
	} else {
		bar = &noopProgressBar{}
	}

	f, err := os.Open(input)
	if err != nil {
		return err
	}
	defer f.Close()

	var header cryptolib.FileHeader
	if err := binary.Read(f, binary.BigEndian, &header); err != nil {
		return err
	}

	decryptor, err := cryptolib.NewDecryptor(password, header.Salt[:])
	if err != nil {
		return err
	}

	if _, err := f.Seek(0, 0); err != nil {
		return err
	}

	reader := &progressReader{
		r:     f,
		bar:   bar,
		total: fileSize,
	}

	outFile, err := os.Create(output)
	if err != nil {
		return err
	}
	defer outFile.Close()

	if err := decryptor.Decrypt(reader, outFile); err != nil {
		_ = bar.Clear()
		return err
	}

	_ = bar.Finish()
	ui.PrintSuccess(output, fileSize)
	return nil
}
