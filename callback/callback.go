package callback

import "github.com/DeathHand/gateway/model"

/**
Callback interface represents call-back functionality
*/
type Callback interface {
	Add(message *model.Message)
	Run()
}
