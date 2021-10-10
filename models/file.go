package models

import (
	"encoding/json"

	"github.com/artem-malko/auth-and-go/infrastructure/convert"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type LinkedEntityType string

func (l LinkedEntityType) String() string {
	return string(l)
}

const (
	LinkedEntityTypeAccountAvatar LinkedEntityType = "account_avatar"
)

type FileModerationStatus string

func (f FileModerationStatus) String() string {
	return string(f)
}

const (
	FileModerationStatusRejected   FileModerationStatus = "rejected"
	FileModerationStatusInProgress FileModerationStatus = "in_progress"
	FileModerationStatusApproved   FileModerationStatus = "approved"
	FileModerationStatusNotStarted FileModerationStatus = "not_started"
)

type File struct {
	ID               uuid.UUID            `json:"-"`
	LinkedEntityID   uuid.UUID            `json:"-"`
	LinkedEntityType LinkedEntityType     `json:"-"`
	AuthorID         uuid.UUID            `json:"-"`
	ModerationStatus FileModerationStatus `json:"-"`
	Name             string               `json:"-"`
	Extensions       []string             `json:"-"`
	Description      string               `json:"-"`
	Size             int                  `json:"-"`
	Copyright        string               `json:"-"`
}

func (f *File) IsDefined() bool {
	if f.Name != "" && len(f.Extensions) > 0 {
		return true
	}

	return false
}

func (f *File) GetFileNames() []string {
	res := []string{}

	if !f.IsDefined() {
		return res
	}

	for _, extension := range f.Extensions {
		res = append(res, f.Name+"."+extension)
	}

	return res
}

func (f *File) ConvertToMap(prefix string) map[string]interface{} {
	if !f.IsDefined() {
		return nil
	}

	return map[string]interface{}{
		"id":          f.ID.String(),
		"name":        prefix + f.Name,
		"extensions":  f.Extensions,
		"description": convert.NewStringPointer(f.Description),
		"size":        f.Size,
		"copyright":   convert.NewStringPointer(f.Copyright),
	}
}

func ParseFileFromBytes(rawFile []byte) (*File, error) {
	var parsedFile *File

	err := json.Unmarshal(rawFile, parsedFile)

	if err != nil {
		return nil, errors.Wrap(err, "unmarshal file error")
	}

	return parsedFile, nil
}

func GetAccountAvatarFileFromFiles(files []File) File {
	for _, file := range files {
		if file.LinkedEntityType == LinkedEntityTypeAccountAvatar {
			return file
		}
	}

	return File{}
}
