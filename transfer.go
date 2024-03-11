package main

import (
	"archive/tar"
	"compress/gzip"
	"encoding/binary"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
)

func transfer(dumpType int, ipAddress, checkpointDir string) {
	conn, err := net.Dial("tcp", ipAddress+":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Send the dumpType
	err = binary.Write(conn, binary.BigEndian, int32(dumpType))
	if err != nil {
		log.Fatal("Failed to send dump type:", err)
	}

	// Create a gzip writer
	gzipWriter := gzip.NewWriter(conn)
	defer gzipWriter.Close()

	// Create a tar writer
	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	baseDir := filepath.Dir(checkpointDir)

	// Walk through the checkpoint directory and write each file to the tar archive
	err = filepath.Walk(checkpointDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path == checkpointDir {
			return nil // Skip the root directory
		}

		relPath, err := filepath.Rel(baseDir, path)
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}

		header.Name = filepath.ToSlash(relPath)

		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(tarWriter, file)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	// Ensure all data is flushed to the connection before sending the dumpType
	if err := tarWriter.Close(); err != nil {
		log.Fatal("Failed to close tar writer:", err)
	}
	if err := gzipWriter.Close(); err != nil {
		log.Fatal("Failed to close gzip writer:", err)
	}

	log.Println("Directory and dump type sent successfully")
}

// func main() {
// 	// Example usage
// 	// transfer(1, "172.31.28.114", "/home/ubuntu/tmp/checkpoint")
// 	transfer(1, "172.31.22.243", "/home/ubuntu/mycontainer/checkpoint/")
// }
