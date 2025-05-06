// inventory.go
package inventory

import (
	"github.com/chiyonn/vendiq2/pricer/pkg/spapi"
)

type Inventory struct {
	client spapi.Requester
}

func New(client spapi.Requester) *Inventory {
	return &Inventory{
		client: client,
	}
}
