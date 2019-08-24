package conn

import (
	"fmt"
	"log"
)

func (c *Conn) Log(v ...interface{}) {
	log.Println(fmt.Sprintf("[gateway:%s] ", c.gateway.Name), v)
}
