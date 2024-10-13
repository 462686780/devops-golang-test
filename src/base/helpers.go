package base

import (
	"encoding/binary"
	"encoding/hex"
	"math/rand"
	"time"

	"github.com/yu31/snowflake"
)

var seeder = rand.New(rand.NewSource(time.Now().UnixNano()))

// The private variables
var (
	reqIDWorker *snowflake.Snowflake
)

// NewRequestID for generate a new id worker
func NewRequestID() (int64, string) {
	id, _ := reqIDWorker.Next()

	var buf = make([]byte, 8, 16)
	binary.BigEndian.PutUint64(buf, uint64(id))
	for i := 0; i < 2; i++ {
		buf[i+6] = byte(seeder.Intn(128))
	}
	s := hex.EncodeToString(buf)
	return id, s
}
