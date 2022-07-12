package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrFileDoesntExist       = errors.New("file doesn't exist")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrFileSeeking           = errors.New("error while offsetting file")
)

func Copy(fromPath, toPath string, offset, limit int64) error {

	file, stat, err := readFile(fromPath, offset)
	defer file.Close()

	err = setOffset(offset, file)
	if err != nil {
		return err
	}

	r := bufio.NewReader(file)
	tempBuf := make([]byte, 3)
	totalBuf := make([]byte, 0, getTotalBuffer(limit, stat, offset))

	totalBuf, err = readOriginalFile(r, tempBuf, totalBuf, limit)
	if err != nil {
		return err
	}

	dest, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer dest.Close()

	if limit != 0 {
		_, err = io.CopyN(dest, bufReader{&totalBuf}, limit)
	} else {
		_, err = io.Copy(dest, bufReader{&totalBuf})
	}
	if err != nil {
		return err
	}

	return nil
}

func readOriginalFile(r *bufio.Reader, tempBuf []byte, totalBuf []byte, limit int64) ([]byte, error) {
	totalRead := 0
	for {
		read, err := r.Read(tempBuf)
		if err == io.EOF {
			fmt.Println("Reached EOF")
			break
		}
		if err != nil {
			fmt.Println("Error while reading file")
			return nil, err
		}

		totalBuf = append(totalBuf, tempBuf...)

		totalRead += read
		if limit != 0 {
			if totalRead > int(limit) {
				fmt.Println("Reached the limit")
				break
			}
		}
	}
	return totalBuf, nil
}

func getTotalBuffer(limit int64, stat os.FileInfo, offset int64) int64 {
	var totalBufSize int64
	if limit > 0 && limit < stat.Size() {
		totalBufSize = limit - offset
	} else {
		totalBufSize = stat.Size()
	}
	return totalBufSize
}

func setOffset(offset int64, file *os.File) error {
	if offset < 0 {
		return nil
	}
	_, err := file.Seek(offset, 0)
	if err != nil {
		return ErrFileSeeking
	}
	return nil
}

func readFile(fromPath string, offset int64) (*os.File, os.FileInfo, error) {
	file, err := os.Open(fromPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, ErrFileDoesntExist
		}
		return nil, nil, ErrUnsupportedFile
	}

	stat, err := validateOffset(err, file, offset)

	return file, stat, nil
}

func validateOffset(err error, file *os.File, offset int64) (os.FileInfo, error) {
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	if offset > stat.Size() {
		return nil, ErrOffsetExceedsFileSize
	}
	return stat, nil
}

type bufReader struct {
	payload *[]byte
}

func (b bufReader) Read(p []byte) (n int, err error) {
	l := copy(p, *b.payload)
	return l, io.EOF
}
