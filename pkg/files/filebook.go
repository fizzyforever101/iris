package files

import (
	"bufio"
	"io/ioutil"
	"os"
	"time"
	"io"
	"path/filepath"
	"fmt"

	"github.com/ipfs/go-cid"
	"github.com/pkg/errors"

	"happystoic/p2pnetwork/pkg/org"
)

type FileMeta struct {
	ExpiredAt time.Time
	Expired   bool

	Available bool
	Path      string

	Rights      []*org.Org
	Severity    Severity
	Description interface{}

	ChunkSize    int32     // size (in bytes) of each chunk
	ChunkCount   int32     // total number of chunks
	// a slice of booleans that indicates which chunk has been downloaded (or is available)
	ChunksStatus []bool    // true if the chunk is available locally, false otherwise
}

func GetFileCid(path string) (*cid.Cid, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = f.Close()
	}()
	reader := bufio.NewReader(f)
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return GetBytesCid(data)
}

// GetBytesCid generates a CID from the given byte slice
// using the CIDv0 format. It uses the SHA-256 hash function
func GetBytesCid(data []byte) (*cid.Cid, error) {
	var builder cid.V0Builder
	c, err := cid.V0Builder.Sum(builder, data)
	return &c, err
}

type FileBook struct {
	files map[cid.Cid]*FileMeta
}

func NewFileBook() *FileBook {
	return &FileBook{files: make(map[cid.Cid]*FileMeta)}
}

func (fb *FileBook) Get(cid *cid.Cid) *FileMeta {
	if meta, exists := fb.files[*cid]; exists {
		return meta
	}
	return nil
}

func (fb *FileBook) AddFile(cid *cid.Cid, meta *FileMeta) error {
	if _, exists := fb.files[*cid]; exists {
		return errors.Errorf("file with cid %s already exists", cid.String())
	}
	fb.files[*cid] = meta
	return nil
}

// update whether chunk is available or not
func (fb *FileBook) UpdateChunkStatus(cid *cid.Cid, chunkIndex int32, available bool) error {
	meta := fb.Get(cid)
	if meta == nil {
		return errors.Errorf("file with cid %s not found", cid.String())
	}
	if int(chunkIndex) >= len(meta.ChunksStatus) {
		return errors.Errorf("chunk index %d out of range", chunkIndex)
	}
	meta.ChunksStatus[chunkIndex] = available
	return nil
}

// use to check if all chunks are available to reassemble the file
func (meta *FileMeta) IsComplete() bool {
    for _, available := range meta.ChunksStatus {
        if !available {
            return false
        }
    }
    return true
}

// splits a file into chunks and stores them in the chunk storage
func (fb *FileBook) SplitFileIntoChunks(path string, chunkSize int32, baseDir string) (*cid.Cid, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fileCid, err := GetFileCid(path)
	if err != nil {
		return nil, err
	}

	meta := &FileMeta{
		Path:        path,
		ChunkSize:   chunkSize,
		ChunkCount:  0,
		ChunksStatus: []bool{},
	}

	// Use the provided baseDir for chunk storage
	chunkDir := filepath.Join(baseDir, "chunks", fileCid.String())
	err = os.MkdirAll(chunkDir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, chunkSize)
	for {
		n, err := f.Read(buf)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if n == 0 {
			break
		}

		// Define chunkPath before using it
		chunkPath := filepath.Join(chunkDir, fmt.Sprintf("%d", meta.ChunkCount))
		err = ioutil.WriteFile(chunkPath, buf[:n], os.ModePerm)
		if err != nil {
			return nil, err
		}

		meta.ChunkCount++
		meta.ChunksStatus = append(meta.ChunksStatus, true)
	}

	err = fb.AddFile(fileCid, meta)
	if err != nil {
		return nil, err
	}

	return fileCid, nil
}

// reassembles a file from its chunks and writes it to the output path string)
func (fb *FileBook) ReassembleFile(cid *cid.Cid, outputPath string, baseDir string) error {
	meta := fb.Get(cid)
	if meta == nil {
		return errors.Errorf("file with cid %s not found", cid.String())
	}

	// Ensure all chunks are available
	if (!meta.IsComplete()) {
		return errors.Errorf("not all chunks are available for file %s", cid.String())
	}

	// Open the output file for writing
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// Read and concatenate chunks in orderString())
	chunkDir := filepath.Join(baseDir, "chunks", cid.String())
	for i := int32(0); i < meta.ChunkCount; i++ {
		chunkPath := filepath.Join(chunkDir, fmt.Sprintf("%d", i))
		chunkData, err := os.ReadFile(chunkPath)
		_, err = outputFile.Write(chunkData)
		if err != nil {
			return errors.Errorf("error writing chunk %d to output file for file %s: %s", i, cid.String(), err)
		}
	}

	// Update metadata to mark the file as reassembled
	meta.Available = true
	meta.Path = outputPath
	return nil
}
