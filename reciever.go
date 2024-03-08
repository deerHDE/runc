package main

import (
	"archive/tar"
	"compress/gzip"
	"encoding/binary"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	conn, err := listener.Accept()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Read the dumpType value
	var dumpType int32
	err = binary.Read(conn, binary.BigEndian, &dumpType)
	if err != nil {
		log.Fatal("Failed to read dump type:", err)
	}

	// Receive and unzip the file
	gzipReader, err := gzip.NewReader(conn)
	if err != nil {
		log.Fatal(err)
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			log.Fatal(err)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(header.Name, 0755); err != nil {
				log.Fatal(err)
			}
		case tar.TypeReg:
			file, err := os.Create(header.Name)
			if err != nil {
				log.Fatal(err)
			}
			if _, err := io.Copy(file, tarReader); err != nil {
				log.Fatal(err)
			}
			file.Close()
		}
	}

	log.Printf("File received successfully. Dump type: %d\n", dumpType)

	// Decide what to do with the file based on the dumpType value
	// For example:
	switch dumpType {
	case 1:
		log.Println("Performing action for dump type 1")
		// Add your logic here
	case 2:
		log.Println("Performing action for dump type 2")
	}
}
