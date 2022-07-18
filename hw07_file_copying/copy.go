package main

import (
	"bufio"
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrFileDoesntExist       = errors.New("file doesn't exist")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrFileSeeking           = errors.New("error while offsetting file")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	file, err := readFile(fromPath)
	if err != nil {
		return err
	}
	defer file.Close()

	err = validateOffset(file, offset)
	if err != nil {
		return err
	}
	err = setOffset(offset, file)
	if err != nil {
		return err
	}

	reader := bufio.NewReader(file)
	dest, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer dest.Close()

	var bar *pb.ProgressBar
	stat, _ := file.Stat()
	sizeToCopy := stat.Size()

	if limit > 0 && sizeToCopy > limit {
		sizeToCopy = limit
	}

	bar = pb.Full.Start64(sizeToCopy)
	barReader := bar.NewProxyReader(reader)
	_, err = io.CopyN(dest, barReader, sizeToCopy)

	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}
	bar.Finish()

	return nil
}

func setOffset(offset int64, file *os.File) error {
	if offset < 0 {
		return nil
	}
	if _, err := file.Seek(offset, 0); err != nil {
		return ErrFileSeeking
	}
	return nil
}

func readFile(fromPath string) (*os.File, error) {
	file, err := os.Open(fromPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrFileDoesntExist
		}
		return nil, ErrUnsupportedFile
	}

	return file, nil
}

func validateOffset(file *os.File, offset int64) error {
	stat, err := file.Stat()
	if err != nil {
		return err
	}

	if offset > stat.Size() {
		return ErrOffsetExceedsFileSize
	}
	return nil
}
