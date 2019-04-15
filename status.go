package arigo

import (
	"encoding/json"
	"time"
)

// DownloadStatus represents the status of a download.
type DownloadStatus string

const (
	// StatusActive represents currently downloading/seeding downloads
	StatusActive DownloadStatus = "active"
	// StatusWaiting represents downloads in the queue
	StatusWaiting DownloadStatus = "waiting"
	// StatusPaused represents paused downloads
	StatusPaused DownloadStatus = "paused"
	// StatusError represents downloads that were stopped because of error
	StatusError DownloadStatus = "error"
	// StatusComplete represents stopped and completed downloads
	StatusCompleted DownloadStatus = "completed"
	// StatusRemoved represents the downloads removed by user
	StatusRemoved DownloadStatus = "removed"
)

// Status holds information for a download.
type Status struct {
	GID             string         // gid of the download
	Status          DownloadStatus // Download status
	TotalLength     uint           `json:",string"` // Total length of the download in bytes
	CompletedLength uint           `json:",string"` // Completed length of the download in bytes
	UploadLength    uint           `json:",string"` // Uploaded length of the download in bytes

	// Hexadecimal representation of the download progress.
	// The highest bit corresponds to the piece at index 0. Any set bits indicate loaded pieces,
	// while unset bits indicate not yet loaded and/or missing pieces.
	// Any overflow bits at the end are set to zero.
	// When the download was not started yet, this will be an empty string.
	BitField      string
	DownloadSpeed uint       `json:",string"` // Download speed of this download measured in bytes/sec
	UploadSpeed   uint       `json:",string"` // Upload speed of this download measured in bytes/sec
	InfoHash      string     // InfoHash. BitTorrent only
	NumSeeders    uint       `json:",string"` // The number of seeders aria2 has connected to. BitTorrent only
	Seeder        bool       `json:",string"` // true if the local endpoint is a seeder. Otherwise false. BitTorrent only
	PieceLength   uint       `json:",string"` // Piece length in bytes
	NumPieces     uint       `json:",string"` // The number of pieces
	Connections   uint       `json:",string"` // The number of peers/servers aria2 has connected to
	ErrorCode     ExitStatus `json:",string"` // The code of the last error for this item, if any.
	ErrorMessage  string     // The human readable error message associated to ErrorCode

	// List of GIDs which are generated as the result of this download.
	// For example, when aria2 downloads a Metalink file, it generates downloads described in the Metalink
	// (see the --follow-metalink option). This value is useful to track auto-generated downloads.
	// If there are no such downloads, this will be an empty slice
	FollowedBy []string

	// The reverse link for followedBy.
	// A download included in followedBy has this object’s GID in its following value
	Following string

	// GID of a parent download. Some downloads are a part of another download.
	// For example, if a file in a Metalink has BitTorrent resources,
	// the downloads of “.torrent” files are parts of that parent.
	// If this download has no parent, this will be an empty string
	BelongsTo  string
	Dir        string           // Directory to save files
	Files      []File           // Slice of files.
	BitTorrent BitTorrentStatus // Information retrieved from the .torrent (file). BitTorrent only

	// The number of verified number of bytes while the files are being hash checked.
	// This key exists only when this download is being hash checked
	VerifiedLength         uint `json:",string"`
	VerifyIntegrityPending bool `json:",string"` // true if this download is waiting for the hash check in a queue.
}

// UNIXTime is just time.Time but it marshals to a Unix timestamp.
type UNIXTime struct {
	time.Time
}

func (t UNIXTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Unix())
}

func (t *UNIXTime) UnmarshalJSON(data []byte) error {
	var ts int64
	err := json.Unmarshal(data, &ts)
	if err != nil {
		return err
	}

	*t = UNIXTime{time.Unix(ts, 0)}
	return nil
}

// TorrentMode represents the file mode of the torrent
type TorrentMode string

const (
	// TorrentModeSingle represents the file mode single
	TorrentModeSingle TorrentMode = "single"
	// TorrentModeMulti represents the file mode multi
	TorrentModeMulti TorrentMode = "multi"
)

// BitTorrentStatus holds information for a BitTorrent download
type BitTorrentStatus struct {
	// List of lists of announce URIs.
	// If the torrent contains announce and no announce-list, announce is converted to the announce-list format
	AnnounceList     []URI
	Comment          string               // The comment of the torrent
	CreationDateUNIX UNIXTime             `json:",string"` // The creation time of the torrent
	Mode             TorrentMode          // File mode of the torrent
	Info             BitTorrentStatusInfo // Information from the info dictionary
}

// A BitTorrentStatusInfo holds information from the info dictionary.
type BitTorrentStatusInfo struct {
	Name string // name in info dictionary
}
