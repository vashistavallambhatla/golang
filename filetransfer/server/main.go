package main

import (
	pb "example/hello/filetransfer/grpc"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"path/filepath"
)

type server struct {
	pb.UnimplementedFileTransferServiceServer
	mu sync.Mutex
}

func (s *server) UploadFile(stream pb.FileTransferService_UploadFileServer) error {
	log.Println("Receiving file...")

	firstChunk, err := stream.Recv()
	if err == io.EOF {
		log.Println("File received completely.")
		return stream.SendAndClose(&pb.UploadStatus{
			Success: true,
			Message: "File uploaded successfully",
		})
	}
	if err != nil {
		log.Printf("error receiving chunk: %w", err)
		return err
	}

	fileName := firstChunk.FileName
	fileName = fmt.Sprintf("server/%s.txt", fileName)
	file, err := os.Create(fileName)
	if err != nil {
		log.Printf("failed to create output file: %v", err)
		return err
	}

	defer file.Close()
	
	_, err = file.Write(firstChunk.ChunkData)
	if err != nil {
		log.Printf("error writing to the file %v", err)
		return err
	}

	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			log.Printf("successfully written to the file %v", fileName)
			return stream.SendAndClose(
				&pb.UploadStatus{
					Success: true,
					Message: "Successfully uploaded the file",
				})
		}
		if err != nil {
			log.Printf("error receiving the data: %v",err)
			return err
		}

		_, err = file.Write(chunk.ChunkData)
		if err != nil {
			log.Printf("error uploading the file: %v",err)
			return err
		}
		log.Printf("received chunk %d",chunk.ChunkData)
	}
}

func (s * server) DownloadFile(downloadReq *pb.FileRequest,stream pb.FileTransferService_DownloadFileServer) error  {
	fileName := filepath.Base(downloadReq.FileName)
	absPath, err := filepath.Abs("server")
	if err != nil {
		log.Print()
	}

	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("no such file found %v",err)
		return err
	}
	defer file.Close()

	chunkSize := downloadReq.ChunkSize
	buf := make([]byte,chunkSize)
	var chunkIndex int32

	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			log.Print("Successfully sent the file")
			return stream.Send(&pb.FileChunk{
				FileName: fileName,
				ChunkData: []byte{}, // apparently this is more protobuf safe than sending a nil
				IsLastChunk: true,
				ChunkIndex: chunkIndex,
			})
		}
		if err != nil {
			log.Printf("error reading from the file %v",fileName)
			return err
		}
		
		if sendErr := stream.Send(
			&pb.FileChunk{
				FileName: fileName,
				ChunkData: buf[:n],
				IsLastChunk: false,
				ChunkIndex: chunkIndex,
			},
		); sendErr != nil {
			log.Printf("Error sending chunk: %v",sendErr)
			return sendErr
		}
		chunkIndex++
	}
} 

func main() {

}
