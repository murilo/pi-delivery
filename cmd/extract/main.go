package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/googlecloudplatform/pi-delivery/gen/index"
	"github.com/googlecloudplatform/pi-delivery/pkg/obj/gcs"
	"github.com/googlecloudplatform/pi-delivery/pkg/unpack"
)

func main() {
	start := flag.Int64("s", 0, "Start offset")
	n := flag.Int64("n", 100, "Number of digits to read")
	flag.Parse()

	if *n <= 0 {
		return
	}
	if *start < 0 {
		*start += index.Decimal.TotalDigits()
	}

	ctx := context.Background()
	sc, err := gcs.NewClient(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't initialize storage client: %v\n", err)
		os.Exit(1)
	}
	defer sc.Close()

	unpackReader := unpack.NewReader(ctx, index.Decimal.NewReader(ctx, sc.Bucket(index.BucketName)))

	var reader io.Reader
	if _, err := unpackReader.Seek(*start, io.SeekStart); err != nil {
		fmt.Fprintf(os.Stderr, "seek failed: %v\n", err)
		os.Exit(1)
	}
	reader = unpackReader

	blocoraw := new(bytes.Buffer)
	written, err := io.CopyN(blocoraw, reader, *n)
	if err != nil {
		fmt.Fprintf(os.Stderr, "I/O error: %v\n", err)
		os.Exit(1)
	}
	var pos int64 = 0
	bloco := blocoraw.String()
	for pos = 0; pos < written-22; pos++ {
		ispal := true
		for i := 0; i < 21; i++ {
			if bloco[pos+int64(i)] != bloco[pos+21-int64(i)-1] {
				ispal = false
				break
			}
		}
		if ispal {
			println("palindrome found at", *start+pos)
		}
		ispal = true
		for i := 0; i < 22; i++ {
			if bloco[pos+int64(i)] != bloco[pos+22-int64(i)-1] {
				ispal = false
				break
			}
		}
		if ispal {
			println("palindrome found at", *start+pos)
		}
	}
}
