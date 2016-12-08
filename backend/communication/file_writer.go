package communication

import (
	"encoding/binary"
	"log"
	"os"
	"path"
	"superstellar/backend/events"
	"superstellar/backend/pb"
	"superstellar/backend/state"
	"time"

	"github.com/golang/protobuf/proto"
)

const (
	DIRECTORY = "log"
)

type FileWriter struct {
	file  *os.File
	space *state.Space
	ch    chan *pb.Space
}

func NewFileWriter(space *state.Space) (*FileWriter, error) {
	file, err := openFile()

	if err != nil {
		return nil, err
	}

	return &FileWriter{
		file:  file,
		space: space,
		ch:    make(chan *pb.Space, 100),
	}, nil
}

func openFile() (*os.File, error) {
	filename := time.Now().Format("2006-01-02_150405.log")

	if err := os.MkdirAll(DIRECTORY, 0777); err != nil {
		log.Fatal(err)
		return nil, err
	}

	filepath := path.Join(DIRECTORY, filename)
	file, err := os.Create(filepath)

	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("log file created:", filepath)
	}

	return file, err
}

func (fileWriter *FileWriter) Run() {
	defer fileWriter.close()

	for {
		protoSpace := <-fileWriter.ch
		fileWriter.writeToFile(protoSpace)
	}
}

func (fileWriter *FileWriter) writeToFile(space *pb.Space) error {
	data, err := proto.Marshal(space)
	if err != nil {
		log.Fatal(err)
		return err
	}

	msgLen := uint32(len(data))
	if binary.Write(fileWriter.file, binary.LittleEndian, msgLen); err != nil {
		log.Fatal(err)
		return err
	}

	if fileWriter.file.Write(data); err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func (fileWriter *FileWriter) close() {
	fileWriter.file.Close()
}

func (fileWriter *FileWriter) HandlePhysicsReady(physicsReadyEvent *events.PhysicsReady) {
	fileWriter.ch <- fileWriter.space.ToProto(true)
}

func (fileWriter *FileWriter) HandleTimeTick(timeTickEvent *events.TimeTick) {
	fileWriter.ch <- fileWriter.space.ToProto(true)
}
