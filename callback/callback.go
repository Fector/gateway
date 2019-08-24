package callback

import "github.com/DeathHand/gateway/model"

/**
Callback interface represents call-back functionality
*/
type Callback interface {
	Send(message *model.Message)
	Error() *chan error
}
