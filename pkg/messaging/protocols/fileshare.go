package protocols

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/pkg/errors"

	"happystoic/p2pnetwork/pkg/config"
	ldht "happystoic/p2pnetwork/pkg/dht"
	"happystoic/p2pnetwork/pkg/files"
	"happystoic/p2pnetwork/pkg/messaging/pb"
	"happystoic/p2pnetwork/pkg/messaging/utils"
	"happystoic/p2pnetwork/pkg/org"
)

/*
Implementation TODOs:
- check completion: update FileMeta data to mark downloaded chunks
- file reassembly (optional): if all chunks are available for a file, reassemble to stoe complete file
*/

// p2p protocol definition
const p2pFileShareMetadataProtocol = "/fileShare-metadata/0.0.1"
const p2pFileShareDownloadProtocol = "/fileShare-download/0.0.1"

// FileShareProtocol type
type FileShareProtocol struct {
	*utils.ProtoUtils

	downloadDir string

	fileBook *files.FileBook
	dht      *ldht.Dht
	spreader *Spreader
}

type Tl2NlRedisFileShareAnnounce struct {
    ExpiredAt       int64       `json:"expired_at"`
    Description     interface{} `json:"description"`
    Severity        string      `json:"severity"`
    Path            string      `json:"path"`
    Rights          []string    `json:"rights"`
    TotalSize       int64       `json:"total_size"`      // total file size in bytes
    ChunkSize       int32       `json:"chunk_size"`      // size of each chunk
    ChunkCount      int32       `json:"chunk_count"`     // total number of chunks
    AvailableChunks []int32     `json:"available_chunks"`// indices of chunks available
}

type Tl2NlRedisFileShareDownloadReq struct {
    FileId       string   `json:"file_id"`
    ChunkIndices []int32  `json:"chunk_indices"` // list of chunk indices requested
}

type Nl2TlRedisFileShareMetadata struct {
    FileId          string             `json:"file_id"`
    Severity        string             `json:"severity"`
    Sender          utils.PeerMetadata `json:"sender"`
    Description     interface{}        `json:"description"`
    TotalSize       int64              `json:"total_size"`      // total file size in bytes
    ChunkSize       int32              `json:"chunk_size"`      // size of each chunk
    ChunkCount      int32              `json:"chunk_count"`     // total number of chunks
    AvailableChunks []int32            `json:"available_chunks"`// indices of chunks available
}

type Nl2TlRedisFileShareDownloadDone struct {
    FileId     string             `json:"file_id"`
    Path       string             `json:"path"`
    Sender     utils.PeerMetadata `json:"sender"`
    ChunkIndex int32              `json:"chunk_index,omitempty"` // indicate which chunk was downloaded
}

func NewFileShareProtocol(ctx context.Context, pu *utils.ProtoUtils, fb *files.FileBook,
	dht *ldht.Dht, cfg *config.FileShareSettings) *FileShareProtocol {

	spreader := NewSpreader(ctx, pu, cfg.MetaSpreadSettings)

	fs := &FileShareProtocol{pu, cfg.DownloadDir, fb, dht, spreader}

	_ = fs.RedisClient.SubscribeCallback("tl2nl_file_share", fs.onRedisFileAnnouncement)
	_ = fs.RedisClient.SubscribeCallback("tl2nl_file_share_download", fs.onDownloadRequest)
	fs.Host.SetStreamHandler(p2pFileShareMetadataProtocol, fs.onP2PMetadata)
	fs.Host.SetStreamHandler(p2pFileShareDownloadProtocol, fs.onP2PDownload)
	return fs
}

func (fs *FileShareProtocol) onDownloadRequest(data []byte) {
	fileAnnouncement := Tl2NlRedisFileShareDownloadReq{}
	err := json.Unmarshal(data, &fileAnnouncement)
	if err != nil {
		log.Errorf("error unmarshalling Tl2NlRedisFileShareDownloadReq from redis: %s", err)
		return
	}
	log.Debug("received file download request message from TL")

	fileCid, err := cid.Decode(fileAnnouncement.FileId)
	if err != nil {
		log.Errorf("error decoding file cid: %s", err)
		return
	}
	meta := fs.fileBook.Get(&fileCid)
	if meta == nil {
		log.Errorf("file with cid %s has no stored metadata", fileCid.String())
		return
	}
	if meta.Available && meta.Path != "" {
		log.Errorf("file with cid %s is already available locally %s", fileCid.String(), meta.Path)
		return
	}
	// TODO shall I also check if I have rights for the file? Or can I assume that?
	providers, err := fs.dht.GetProvidersOf(fileCid)
	if err != nil {
		log.Errorf("error getting providers of file %s: %s", fileCid.String(), err)
		return
	}
	if len(providers) == 0 {
		log.Errorf("Found no providers of %s in DHT", fileCid.String())
		return
	}
	// sort providers based on their reliability to decreasing order
	fs.ReliabilitySort(providers)

	// if specific chunk indices are requested, handle chunked download
	if len(fileAnnouncement.ChunkIndices) > 0 {
        // for each requested chunk, initiate a chunk download
        for _, chunkIndex := range fileAnnouncement.ChunkIndices {
            path, sender := fs.downloadSingleChunk(fileCid, chunkIndex)
            if path == "" || sender == nil {
                log.Errorf("failed to download chunk %d of file %s", chunkIndex, fileCid.String())
                continue
            }
            // notify TL about the downloaded chunk
            err = fs.notifyTLAboutDownload(fileCid, *sender, path, chunkIndex)
            if err != nil {
                log.Errorf("error sending download confirmation to redis for chunk %d: %s", chunkIndex, err)
                continue
            }
            log.Infof("successfully downloaded chunk %d of file %s", chunkIndex, fileCid.String())
        }
    } else {
        // otherwise, proceed with the full file download as in original implementation
		// now use DHT to download the file
        path, sender := fs.downloadFile(providers, fileCid)
        if path == "" || sender == nil {
			// did not succeed
            return
        }
		// notify TL where file is downloaded
        err = fs.notifyTLAboutDownload(fileCid, *sender, path, -1)
        if err != nil {
            log.Errorf("error sending download confirmation to redis: %s", err)
            return
        }
        meta.Available = true
        meta.Path = path
        err = fs.dht.StartProviding(fileCid)
        if err != nil {
            log.Errorf("error starting providing file in DHT: %s", err)
        }
        log.Infof("successfully downloaded the file %s to path %s", fileCid.String(), path)
    }

	// Check if all chunks are downloaded
	meta = fs.fileBook.Get(&fileCid)
	if meta != nil && meta.IsComplete() {
		outputPath := fmt.Sprintf("%s/%s", fs.downloadDir, fileCid.String())
		err := fs.fileBook.ReassembleFile(&fileCid, outputPath, fs.downloadDir)
		if err != nil {
			log.Errorf("error reassembling file %s: %s", fileCid.String(), err)
			return
		}
		log.Infof("successfully reassembled file %s to path %s", fileCid.String(), outputPath)
		meta.Available = true
		meta.Path = outputPath
	}
}

// helper function to download individual chunk
func (fs *FileShareProtocol) downloadSingleChunk(fileCid cid.Cid, chunkIndex int32) (string, *utils.PeerMetadata) {
    // create a FileDownloadRequest specifically for this chunk
    req := &pb.FileDownloadRequest{
        Cid:          fileCid.String(),
        ChunkSize:    fs.getChunkSizeFor(fileCid), // Define this helper to retrieve the chunk size from metadata
        ChunkIndices: []int32{chunkIndex},
    }
    // Sabina's TODO: adapt exisiting P2P download logic from tryFileProvider
    path, err := fs.tryDownloadChunk(req, fileCid) 
    if err != nil {
        log.Errorf("error downloading chunk %d: %s", chunkIndex, err)
        return "", nil
    }
    // return the local path where the chunk is stored and the peer metadata
    return path, fs.getLastSuccessfulProvider()
}

// not sure if this is necessary/redundant because FileBook already has a getter
func (fs *FileShareProtocol) getChunkSizeFor(fileCid cid.Cid) int32 {
    meta := fs.fileBook.Get(&fileCid)
    if meta == nil {
        log.Errorf("metadata not found for file cid: %s", fileCid.String())
        return 0
    }
    return meta.ChunkSize
}

// helper to download individual chunks by calling tryFileProvider
func (fs *FileShareProtocol) tryDownloadChunk(req *pb.FileDownloadRequest, fileCid cid.Cid) (string, error) {
    // get providers from the DHT for this file
    providers, err := fs.dht.GetProvidersOf(fileCid)
    if err != nil {
        return "", errors.Errorf("error getting providers of file %s: %s", fileCid.String(), err)
    }
    if len(providers) == 0 {
        return "", errors.Errorf("no providers found for file %s", fileCid.String())
    }

    // try each provider until one returns the chunk successfully
    for _, p := range providers {
        // use tryFileProvider to send the request/get the response
        path, err := fs.tryFileProvider(req, p.ID, fileCid)
        if err != nil {
            log.Error(err)
            continue
        }
        
        // If we successfully downloaded the chunk, update the ChunksStatus in the metadata
        if len(req.ChunkIndices) > 0 {
            meta := fs.fileBook.Get(&fileCid)
            if meta != nil {
                for _, chunkIndex := range req.ChunkIndices {
                    err = fs.fileBook.UpdateChunkStatus(&fileCid, chunkIndex, true)
                    if err != nil {
                        log.Errorf("error updating chunk status for chunk %d: %s", chunkIndex, err)
                        // Continue even if there's an error updating the status
                    } else {
                        log.Debugf("updated chunk status for chunk %d of file %s", chunkIndex, fileCid.String())
                    }
                }
            } else {
                log.Errorf("metadata not found for file %s after download", fileCid.String())
            }
        }
        
        // If we get a valid path (i.e. the chunk was downloaded), return it.
        return path, nil
    }
    return "", errors.Errorf("failed to download chunk from all providers for file %s", fileCid.String())
}

func (fs *FileShareProtocol) createP2PFileDownloadReq(fileCid cid.Cid) (*pb.FileDownloadRequest, error) {
	msgMetaData, err := fs.NewProtoMetaData()
	if err != nil {
		return nil, errors.WithMessage(err, "error generating new proto metadata: ")
	}

	protoMsg := &pb.FileDownloadRequest{
		Metadata: msgMetaData,
		Cid:      fileCid.String(),
	}
	signature, err := fs.SignProtoMessage(protoMsg)
	if err != nil {
		return nil, errors.WithMessage(err, "error generating signature for new alert message: ")
	}
	protoMsg.Metadata.Signature = signature
	return protoMsg, err
}

func (fs *FileShareProtocol) tryFileProvider(msg proto.Message, p peer.ID, fileCid cid.Cid) (string, error) {
	// TODO disconnect from him after he fails?
	log.Debugf("trying to download %s from %s", fileCid.String(), p.String())

	s, err := fs.InitiateStream(p, p2pFileShareDownloadProtocol, msg)
	if err != nil {
		return "", errors.Errorf("error sending req to %s: %s", p.String(), err)
	}
	_ = s.CloseWrite()
	defer s.Close()

	resp := &pb.FileDownloadResponse{}
	err = fs.DeserializeMessageFromStream(s, resp, false)
	if err != nil {
		return "", errors.Errorf("error reading response: %s", err)
	}

	err = fs.AuthenticateMessage(resp, resp.Metadata)
	if err != nil {
		return "", errors.Errorf("error authenticating response: %s", err)
	}

	if len(resp.Data) == 0 {
		return "", errors.Errorf("peer %s responded %s without data,", p.String(), resp.Status)
	}

	// check if the hash (cid) actually matches
	receivedCid, err := files.GetBytesCid(resp.Data)
	if !receivedCid.Equals(fileCid) {
		err = fs.ReportPeer(p, "provided file with not matching hash")
		if err != nil {
			log.Errorf("error reporting peer: %s", err)
		}
		return "", errors.Errorf("peer %s provided not matching file!", p.String())
	}

	return fs.writeFile(fileCid, resp.Data)
}

func (fs *FileShareProtocol) writeFile(fileCid cid.Cid, data []byte) (string, error) {
	path := fmt.Sprintf("%s/%s", fs.downloadDir, fileCid.String())
	err := os.WriteFile(path, data, 0644)
	if err != nil {
		return "", err
	}
	return path, nil
}

func (fs *FileShareProtocol) downloadFile(providers []peer.AddrInfo, fileCid cid.Cid) (string, *peer.ID) {
	// TODO maybe I have to firstly run Connect(peer)? or maybe at least put addrInfo to peer book?
	// 	    So it does not have to be found in DHTs
	reqMsg, err := fs.createP2PFileDownloadReq(fileCid)
	if err != nil {
		log.Errorf("error generationg file download req: %s", err)
		return "", nil
	}

	for _, p := range providers {
		path, err := fs.tryFileProvider(reqMsg, p.ID, fileCid)
		if err != nil {
			log.Error(err)
			continue
		}
		return path, &p.ID
	}
	return "", nil
}

func (fs *FileShareProtocol) notifyTLAboutDownload(cid cid.Cid, sender peer.ID, path string, chunkIndex int32) error {
	msg := Nl2TlRedisFileShareDownloadDone{
		FileId:     cid.String(),
		Path:       path,
		Sender:     fs.MetadataOfPeer(sender),
		ChunkIndex: chunkIndex, // Include chunk index
	}
	channel := "nl2tl_file_share_downloaded"
	return fs.RedisClient.PublishMessage(channel, msg)
}

func (fs *FileShareProtocol) onP2PDownload(s network.Stream) {
	defer s.Close()
	remote := s.Conn().RemotePeer()
	log.Infof("opened p2p file download request from %s", remote.String())

	req := &pb.FileDownloadRequest{}

	err := fs.DeserializeMessageFromStream(s, req, false)
	if err != nil {
		log.Errorf("error deserilising download req message from stream: %s", err)
		return
	}

	err = fs.AuthenticateMessage(req, req.Metadata)
	if err != nil {
		log.Errorf("error authenticating alert message: %s", err)
		return
	}

	fileCid, err := cid.Decode(req.Cid)
	if err != nil {
		log.Errorf("error decoding cid: %s", err)
		return
	}
	meta := fs.fileBook.Get(&fileCid)
	if meta == nil {
		log.Errorf("unknown cid %s", req.Cid)
		return
	}

	// Validate requested chunk indices
	for _, chunkIndex := range req.ChunkIndices {
		if chunkIndex < 0 || chunkIndex >= meta.ChunkCount {
			log.Errorf("invalid chunk index %d requested for file %s", chunkIndex, req.Cid)
			return
		}
		if !meta.ChunksStatus[chunkIndex] {
			log.Errorf("chunk %d of file %s is not available", chunkIndex, req.Cid)
			return
		}
	}

	// Send each requested chunk
	for _, chunkIndex := range req.ChunkIndices {
		chunkPath := filepath.Join(fs.downloadDir, "chunks", fileCid.String(), fmt.Sprintf("%d", chunkIndex))
		chunkData, err := os.ReadFile(chunkPath)
		if err != nil {
			log.Errorf("error reading chunk %d for file %s: %s", chunkIndex, req.Cid, err)
			return
		}

		resp := &pb.FileDownloadResponse{
			Metadata:   req.Metadata,
			Status:     "OK",
			Data:       chunkData,
			ChunkIndex: chunkIndex,
		}

		err = fs.WriteProtoMsg(resp, s)
		if err != nil {
			log.Errorf("error sending chunk %d of file %s: %s", chunkIndex, req.Cid, err)
			return
		}
		log.Infof("successfully sent chunk %d of file %s", chunkIndex, req.Cid)
	}
}

func (fs *FileShareProtocol) createFileDownloadResp(status string, path string) (*pb.FileDownloadResponse, error) {
	msgMetaData, err := fs.NewProtoMetaData()
	if err != nil {
		return nil, errors.WithMessage(err, "error generating new proto metadata: ")
	}
	data := make([]byte, 0)
	if path != "" {
		data, err = os.ReadFile(path)
		if err != nil {
			return nil, err
		}
	}
	resp := &pb.FileDownloadResponse{
		Metadata: msgMetaData,
		Status:   status,
		Data:     data,
	}

	signature, err := fs.SignProtoMessage(resp)
	if err != nil {
		return nil, errors.WithMessage(err, "error generating signature for new file share resp: ")
	}
	resp.Metadata.Signature = signature
	return resp, nil
}

func (fs *FileShareProtocol) onP2PMetadata(s network.Stream) {
	log.Infof("received p2p file metadata message")
	p2pMeta := &pb.FileMetadata{}

	// TODO following lines are the same as in alert. dont duplicate code
	err := fs.DeserializeMessageFromStream(s, p2pMeta, true)
	if err != nil {
		log.Errorf("error deserilising file metadata proto message from stream: %s", err)
		return
	}

	if fs.WasMsgSeen(p2pMeta.Metadata.Id) {
		log.Debugf("received already seen file metadata message, forwarded by %s", s.Conn().RemotePeer())
		return
	}
	fs.NewMsgSeen(p2pMeta.Metadata.Id, s.Conn().RemotePeer())

	err = fs.AuthenticateMessage(p2pMeta, p2pMeta.Metadata)
	if err != nil {
		log.Errorf("error authenticating file metadata message: %s", err)
		return
	}

	log.Debugf("Received file metadata message authored by %s and "+
		"forwarded by %s", p2pMeta.Metadata.OriginalSender.NodeId, s.Conn().RemotePeer())

	meta, err := fs.fileMetaFromP2P(p2pMeta)
	if err != nil {
		log.Errorf("error creating metadata from p2p msg: %s", err)
		return
	}
	fileCid, err := cid.Decode(p2pMeta.Cid)
	if err != nil {
		log.Errorf("error decoding cid: %s", err)
		return
	}

	err = fs.fileBook.AddFile(&fileCid, meta)
	if err != nil {
		log.Error(err)
		return
	}
	sender, err := peer.Decode(p2pMeta.Metadata.OriginalSender.NodeId)
	if err != nil {
		log.Errorf("error decoding original sender peer id: %s", err)
		return
	}

	err = fs.notifyTLAboutMetadata(fileCid, sender, meta)
	if err != nil {
		log.Errorf("error sending to Redis metadata info: %s", err)
	}

	fs.spreader.startSpreading(p2pFileShareMetadataProtocol, meta.Severity, meta.Rights, p2pMeta, s.Conn().RemotePeer())
	log.Infof("handler onP2PMetadata finished")
}

func (fs *FileShareProtocol) notifyTLAboutMetadata(cid cid.Cid, origSender peer.ID, meta *files.FileMeta) error {
	msg := Nl2TlRedisFileShareMetadata{
		FileId:      cid.String(),
		Severity:    meta.Severity.String(),
		Sender:      fs.MetadataOfPeer(origSender),
		Description: meta.Description,
	}
	channel := "nl2tl_file_share_received_metadata"
	return fs.RedisClient.PublishMessage(channel, msg)
}

func (fs *FileShareProtocol) onRedisFileAnnouncement(data []byte) {
	fileAnnouncement := Tl2NlRedisFileShareAnnounce{}
	err := json.Unmarshal(data, &fileAnnouncement)
	if err != nil {
		log.Errorf("error unmarshalling Tl2NlRedisFileShareAnnounce from redis: %s", err)
		return
	}
	log.Debug("received file share announcement message from TL")

	// Split the file into chunks and prepare metadata
	fileCid, err := fs.fileBook.SplitFileIntoChunks(fileAnnouncement.Path, fileAnnouncement.ChunkSize, fs.downloadDir)
	if err != nil {
		log.Errorf("error splitting file into chunks: %s", err)
		return
	}

	meta := fs.fileBook.Get(fileCid)
	if meta == nil {
		log.Errorf("error retrieving metadata for file %s after splitting", fileCid.String())
		return
	}

	// Update metadata with additional details from the announcement
	meta.ExpiredAt = time.Unix(fileAnnouncement.ExpiredAt, 0)
	meta.Severity, err = files.SeverityFromString(fileAnnouncement.Severity)
	if err != nil {
		log.Errorf("error parsing severity: %s", err)
		return
	}
	meta.Description = fileAnnouncement.Description
	meta.Rights = make([]*org.Org, 0, len(fileAnnouncement.Rights))
	for _, strOrg := range fileAnnouncement.Rights {
		o, err := org.Decode(strOrg)
		if err != nil {
			log.Errorf("error decoding rights: %s", err)
			return
		}
		meta.Rights = append(meta.Rights, o)
	}

	// Start providing the file in the DHT
	err = fs.dht.StartProviding(*fileCid)
	if err != nil {
		log.Errorf("error starting providing file in dht: %s", err)
		retudht	}
	log.Debugf("successfully started providing file %s", fileCid.String())

	// Create and spread P2P metadata
	protoMsg, err := fs.createP2PMeta(*fileCid, fileAnnouncement)
	if err != nil {
		log.Errorf("error creating p2p proto metadata: %s", err)
		return
	}

	// store this msg as seen in casmsgmes back from another peer
	fs.NewMsgSeen(protoMsg.Metadata.Id, fs.Host.ID())

	fs.spreader.startSpreading(p2pFileShareMetadataProtocol, meta.Severity, meta.Rights, protoMsg, fs.Host.ID())
	log.Debugf("handling file share annoucment from TL ended")
}

func (fs *FileShareProtocol) createP2PMeta(fCid cid.Cid, meta Tl2NlRedisFileShareAnnounce) (*pb.FileMetadata, error) {
	msgMetaData, err := fs.NewProtoMetaData()
	if err != nil {
		return nil, errors.WithMessage(err, "error generating new proto metadata: ")
	}

	bytesDesc, err := json.Marshal(meta.Description)
	if err != nil {
		return nil, err
	}

	// add chunk data to protoMsg
	protoMsg := &pb.FileMetadata{
		Metadata:    msgMetaData,
		Cid:         fCid.String(),
		Description: bytesDesc,
		Rights:      meta.Rights,
		Severity:    meta.Severity,
		ExpiredAt:   meta.ExpiredAt,
		TotalSize:	 meta.TotalSize,
		ChunkSize:	 meta.ChunkSize,
		ChunkCount:  meta.ChunkCount,
		AvailableChunks: meta.AvailableChunks,
	}
	signature, err := fs.SignProtoMessage(protoMsg)
	if err != nil {
		return nil, errors.WithMessage(err, "error generating signature for new file share meta: ")
	}
	protoMsg.Metadata.Signature = signature
	return protoMsg, err
}

func (fs *FileShareProtocol) fileMetaFromP2P(p2pMeta *pb.FileMetadata) (*files.FileMeta, error) {
	expiredAt := time.Unix(p2pMeta.ExpiredAt, 0)

	severity, err := files.SeverityFromString(p2pMeta.Severity)
	if err != nil {
		return nil, err
	}

	rights := make([]*org.Org, 0, len(p2pMeta.Rights))
	for _, strOrg := range p2pMeta.Rights {
		o, err := org.Decode(strOrg)
		if err != nil {
			return nil, err
		}
		rights = append(rights, o)
	}

	var desc interface{}
	err = json.Unmarshal(p2pMeta.Description, &desc)
	if err != nil {
		return nil, err
	}

	// Initialize the ChunksStatus slice based on ChunkCount
	chunksStatus := make([]bool, p2pMeta.ChunkCount)
	
	// Mark available chunks as true in ChunksStatus
	for _, chunkIndex := range p2pMeta.AvailableChunks {
		if chunkIndex >= 0 && int(chunkIndex) < len(chunksStatus) {
			chunksStatus[chunkIndex] = true
		}
	}

	meta := &files.FileMeta{
		ExpiredAt:     expiredAt,
		Expired:       time.Now().After(expiredAt),
		Available:     false,
		Path:          "",
		Rights:        rights,
		Severity:      severity,
		Description:   desc,
		ChunkSize:     p2pMeta.ChunkSize,     // Set chunk size from metadata
		ChunkCount:    p2pMeta.ChunkCount,    // Set chunk count from metadata
		ChunksStatus:  chunksStatus,          // Set chunk status from metadata
	}
	return meta, nil
}

func (fs *FileShareProtocol) fileMetaFromRedis(ann *Tl2NlRedisFileShareAnnounce) (*cid.Cid, *files.FileMeta, error) {
	expiredAt := time.Unix(ann.ExpiredAt, 0)
	fileCid, err := files.GetFileCid(ann.Path)
	if err != nil {
		return nil, nil, err
	}

	severity, err := files.SeverityFromString(ann.Severity)
	if err != nil {
		return nil, nil, err
	}

	rights := make([]*org.Org, 0, len(ann.Rights))
	for _, strOrg := range ann.Rights {
		o, err := org.Decode(strOrg)
		if err != nil {
			return nil, nil, err
		}
		rights = append(rights, o)
	}

	// Initialize the ChunksStatus slice based on ChunkCount	var desc interface{}
	chunksStatus := make([]bool, ann.ChunkCount)	err = json.Unmarshal(ann.Description, &desc)
		if err != nil {
	// Mark available chunks as true in ChunksStatus		return nil, nil, err
	for _, chunkIndex := range ann.AvailableChunks {	}
		if chunkIndex >= 0 && int(chunkIndex) < len(chunksStatus) {
			chunksStatus[chunkIndex] = true	// Initialize the ChunksStatus slice based on ChunkCount
		}	chunksStatus := make([]bool, ann.ChunkCount)
	}	
	// Mark available chunks as true in ChunksStatus
	meta := &files.FileMeta{	for _, chunkIndex := range ann.AvailableChunks {
		ExpiredAt:     expiredAt,		if chunkIndex >= 0 && int(chunkIndex) < len(chunksStatus) {
		Expired:       time.Now().After(expiredAt),
		Available:     true,
		Path:          ann.Path,
		Rights:        rights,
		Severity:      severity,
		Description:   ann.Description,
		ChunkSize:     ann.ChunkSize,      // Set chunk size from announcementer(expiredAt),
		ChunkCount:    ann.ChunkCount,     // Set chunk count from announcement
		ChunksStatus:  chunksStatus,       // Set chunk status from announcement
	}
	return fileCid, meta, nilSeverity:      severity,
}
		ChunksStatus:  chunksStatus,
	}
	return fileCid, meta, nil
}
