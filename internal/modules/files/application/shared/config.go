package shared

import (
	"time"

	"github.com/google/uuid"
	"github.com/zchelalo/neuraclinic-file-management/internal/shared/uuidx"
)

type Config struct {
	StorageProvider string
	UploadURLTTL    time.Duration
	DownloadURLTTL  time.Duration
}

type Runtime struct {
	Now     func() time.Time
	NewUUID func() uuid.UUID
}

func DefaultRuntime() Runtime {
	return Runtime{
		Now:     time.Now,
		NewUUID: uuidx.New,
	}
}

func (r Runtime) Normalize() Runtime {
	if r.Now == nil {
		r.Now = time.Now
	}
	if r.NewUUID == nil {
		r.NewUUID = uuidx.New
	}
	return r
}
