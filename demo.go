package main

import (
	"flag"
	"io"
	"log"
	"os"
)

const (
	PROG_NAME       = "simpleEncDec"
	ONE_MB    int64 = 1048576
)

var (
	signByte  int
	maxOpSize int64
)

func init() {
	flag.IntVar(&signByte, "b", 0x38, "Specify the sign byte")
	flag.Int64Var(&maxOpSize, "s", ONE_MB, "Specify the max num of operate size")
}

func main() {
	if len(os.Args) == 1 {
		log.Printf("Usage: %s [options] file-to-be-operated.\n", PROG_NAME)
		return
	}

	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		log.Println("Please specify a file to operate.")
		return
	}

	var xorByte byte
	var realOpSize int64
	var target string = args[0]
	realOpSize = getRealOpSize(target, maxOpSize)

	xorByte = (byte)(signByte)
	log.Println(xorByte, realOpSize)
	operateFile(target, xorByte, realOpSize)
}

func getRealOpSize(name string, max int64) int64 {
	var fSize int64

	fi, err := os.Stat(name)
	if err != nil {
		log.Fatal(err.Error())
	}
	fSize = fi.Size()
	if fSize < max {
		return fSize
	} else {
		return max
	}
}

func operateFile(name string, xor byte, size int64) {
	fp, err := os.OpenFile(name, os.O_RDWR, 0666)
	checkError(err)
	defer fp.Close()

	buff := make([]byte, size)
	num, err := fp.Read(buff)
	if err != nil && err != io.EOF {
		log.Fatal(err.Error())
	}

	for i := 0; i < num; i++ {
		buff[i] ^= xor
	}

	_, err = fp.Seek(0, os.SEEK_SET)
	checkError(err)
	_, err = fp.Write(buff[:num])
	checkError(err)
	log.Printf("Successfully enc/dec, remember the sign is : 0x%x\n", xor)
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
