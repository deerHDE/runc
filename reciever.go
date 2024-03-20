package main

import (
	// "archive/tar"
	// "compress/gzip"
	// "encoding/binary"
	// "io"
	"log"
	"net"
	// "os"
	"encoding/json"
)

// type Message struct {
// 	DumpType int	`json:"dumpType"`
// 	ContainerID string 	`json:"containerID"`
// }


func receive() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for{
		conn, err := listener.Accept()

		if err != nil {
			log.Fatal(err)
			continue
		}
		log.Println("Connection accepted")
		// defer conn.Close()

		// Read the json value
		jsonBytes := make([]byte, 1024)
		n, err := conn.Read(jsonBytes)
		if err != nil {
			log.Fatal(err)
			conn.Close()
			continue
		}

		log.Println("message read from connection")

		var msg Message
		err = json.Unmarshal(jsonBytes[:n], &msg)
		if err != nil {
			log.Fatal(err)
			conn.Close()
			continue
		}


		// var dumpType int32
		// err = binary.Read(conn, binary.BigEndian, &dumpType)
		// if err != nil {
		// 	log.Fatal("Failed to read dump type:", err)
		// }

		// Receive and unzip the file
		// gzipReader, err := gzip.NewReader(conn)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// defer gzipReader.Close()

		// tarReader := tar.NewReader(gzipReader)

		// for {
		// 	header, err := tarReader.Next()
		// 	if err == io.EOF {
		// 		break // End of archive
		// 	}
		// 	if err != nil {
		// 		log.Fatal(err)
		// 	}

		// 	switch header.Typeflag {
		// 	case tar.TypeDir:
		// 		if err := os.Mkdir(header.Name, 0755); err != nil {
		// 			log.Fatal(err)
		// 		}
		// 	case tar.TypeReg:
		// 		file, err := os.Create(header.Name)
		// 		if err != nil {
		// 			log.Fatal(err)
		// 		}
		// 		if _, err := io.Copy(file, tarReader); err != nil {
		// 			log.Fatal(err)
		// 		}
		// 		file.Close()
		// 	}
		// }

		// log.Printf("File received successfully. Dump type: %d\n", dumpType)

		// Decide what to do with the file based on the dumpType value
		// For example:
		dumpType := msg.DumpType
		log.Println("dump type: ", dumpType)

		if dumpType == 1 {
			conn.Close()
			continue
		} else if dumpType == 2 {
			log.Println("received json message with dumptype: ", dumpType)
			conn.Close()
			break
		}
		conn.Close()

	}
	
}
