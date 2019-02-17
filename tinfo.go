package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/zeebo/bencode"
	"gopkg.in/alecthomas/kingpin.v2"
)

type Torrent struct {
	Filepath  string
	Name      string
	Hash      string
	Size      int64
	Trackers  []string
	CreatedOn time.Time
	CreatedBy string
	Comment   string
	Info      Info
}

type Info struct {
	PieceCount int
	PieceSize  int64
	Private    int
	Files      []File
}

type File struct {
	Path string
	Size int64
}

func (torrent *Torrent) Parse(rawTorrent *RawTorrent) {
	info := rawTorrent.Info

	torrent.Name = info.Name
	torrent.CreatedOn = time.Unix(rawTorrent.CreatedOn, 0)
	torrent.CreatedBy = rawTorrent.CreatedBy
	torrent.Comment = rawTorrent.Comment

	torrent.computeHash(rawTorrent)
	torrent.computeSize(rawTorrent)
	torrent.computeTrackers(rawTorrent)

	torrent.Info = Info{
		PieceCount: len(info.Pieces) / 20,
		PieceSize:  info.PieceSize,
		Private:    int(info.Private),
	}

	if len(rawTorrent.Info.Files) == 0 {
		torrent.Info.Files = []File{
			{Path: torrent.Name, Size: torrent.Size},
		}
	} else {
		torrent.Info.Files = []File{}

		for _, file := range info.Files {
			path := filepath.Join(file.Path...)
			path = filepath.Join(torrent.Name, path)
			torrent.Info.Files = append(torrent.Info.Files, File{Path: path, Size: file.Size})
		}
	}
}

func (torrent *Torrent) computeHash(parseTorrent *RawTorrent) {
	hash := sha1.Sum(parseTorrent.RawInfo)
	torrent.Hash = hex.EncodeToString(hash[:])
}

func (torrent *Torrent) computeSize(parseTorrent *RawTorrent) {
	torrent.Size = parseTorrent.Info.Size

	for _, file := range parseTorrent.Info.Files {
		torrent.Size += file.Size
	}
}

func (torrent *Torrent) computeTrackers(parseTorrent *RawTorrent) {
	trackers := map[string]bool{
		parseTorrent.Announce: true,
	}

	for _, group := range parseTorrent.AnnounceList {
		for _, tracker := range group {
			trackers[tracker] = true
		}
	}

	for tracker := range trackers {
		torrent.Trackers = append(torrent.Trackers, tracker)
	}
}

func (torrent *Torrent) ToJSON() (string, error) {
	result, err := json.MarshalIndent(torrent, "", "  ")

	if err != nil {
		return "", err
	}

	return string(result), nil
}

var format = strings.Join([]string{
	"File: %v",
	"Name: %v",
	"Hash: %v",
	"Created By: %v",
	"Created On: %v",
	"Comment: %v",
	"Piece Count: %v",
	"Piece Size: %v",
	"Total Size: %v",
	"Private: %v",
	"Trackers:",
	"%v",
	"Files:",
	"%v",
}, "\n")

var privateFlag = map[int]string{
	0: "no",
	1: "yes",
}

func (torrent *Torrent) ToText() (string, error) {
	var trackers []string

	for _, tracker := range torrent.Trackers {
		trackers = append(trackers, fmt.Sprintf("  %v", tracker))
	}

	var files []string

	for _, file := range torrent.Info.Files {
		files = append(files, fmt.Sprintf("  %v (%v)", file.Path, file.Size))
	}

	return fmt.Sprintf(format,
		torrent.Filepath,
		torrent.Name,
		torrent.Hash,
		torrent.CreatedBy,
		torrent.CreatedOn,
		torrent.Comment,
		torrent.Info.PieceCount,
		torrent.Info.PieceSize,
		torrent.Size,
		privateFlag[torrent.Info.Private],
		strings.Join(trackers, "\n"),
		strings.Join(files, "\n"),
	), nil
}

type RawTorrent struct {
	Announce     string             `bencode:"announce"`
	AnnounceList [][]string         `bencode:"announce-list"`
	Comment      string             `bencode:"comment"`
	CreatedOn    int64              `bencode:"creation date"`
	CreatedBy    string             `bencode:"created by"`
	RawInfo      bencode.RawMessage `bencode:"info"`
	InfoHash     string
	Info         struct {
		PieceSize int64  `bencode:"piece length"`
		Pieces    []byte `bencode:"pieces"`
		Private   int64  `bencode:"private"`
		Name      string `bencode:"name"`
		Size      int64  `bencode:"length"`
		Files     []struct {
			Path []string `bencode:"path"`
			Size int64    `bencode:"length"`
		} `bencode:"files"`
	}
}

func (rt *RawTorrent) Decode(data []byte) error {
	err := bencode.DecodeBytes(data, rt)

	if err != nil {
		return err
	}

	err = bencode.DecodeBytes(rt.RawInfo, &rt.Info)

	if err != nil {
		return err
	}

	if len(rt.Info.Pieces)%20 != 0 {
		return errors.New("malformed torrent")
	}

	return nil
}

type TInfo struct {
	Path     string
	JSONMode bool
}

func (ti TInfo) Run() (string, error) {
	path, err := filepath.Abs(ti.Path)

	if err != nil {
		return "", fmt.Errorf("cannot read the file at %v: %v", path, err)
	}

	data, err := ioutil.ReadFile(path)

	if err != nil {
		return "", fmt.Errorf("cannot read the file at %v: %v", path, err)
	}

	rawTorrent := RawTorrent{}
	err = rawTorrent.Decode(data)

	if err != nil {
		return "", fmt.Errorf("corrupted file at %v: %v", path, err)
	}

	torrent := Torrent{Filepath: path}
	torrent.Parse(&rawTorrent)

	var result string

	if ti.JSONMode {
		result, err = torrent.ToJSON()
	} else {
		result, err = torrent.ToText()
	}

	if err != nil {
		return "", fmt.Errorf("cannot format json: %v", err)
	}

	return result, nil
}

var (
	jsonMode = kingpin.Flag("json", "Enable json output mode.").Short('j').Bool()
	path     = kingpin.Arg("path", "Path to a .torrent file.").Required().String()
)

func main() {
	kingpin.Version("0.0.1")
	kingpin.Parse()

	app := TInfo{
		Path:     *path,
		JSONMode: *jsonMode,
	}

	result, err := app.Run()

	kingpin.FatalIfError(err, "")

	fmt.Println(result)
}
