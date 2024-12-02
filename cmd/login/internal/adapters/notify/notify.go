package notify

import (
	"context"
	"github.com/gofrs/uuid"
)

type Notify struct{}

func (n *Notify) SendNotify(ctx context.Context, userID uuid.UUID) error {
	return nil
}
