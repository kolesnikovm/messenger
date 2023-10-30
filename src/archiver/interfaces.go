package archiver

import "context"

type Archiver interface {
	Start(ctx context.Context)
}
