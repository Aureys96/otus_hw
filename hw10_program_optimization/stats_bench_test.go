package hw10programoptimization

import (
	"archive/zip"
	"log"
	"testing"
)

func BenchmarkGetDomainStat(b *testing.B) {
	reader, err := zip.OpenReader("testdata/users.dat.zip")
	if err != nil {
		log.Fatalln(err)
	}

	data, err := reader.File[0].Open()
	if err != nil {
		log.Fatalln(err)
	}
	defer reader.Close()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = GetDomainStat(data, "com")
	}
}
