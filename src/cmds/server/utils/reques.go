package utils

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/yu31/snowflake"
)

// The private variables
var (
	reqIDWorker *snowflake.Snowflake
)
var seeder = rand.New(rand.NewSource(time.Now().UnixNano()))

func init() {
	// Setup global request id generator worker
	rand.Seed(time.Now().UnixNano())
	num := rand.Intn(100)
	reqIDWorker, _ = snowflake.New(int64(num))
}

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

func Request(ctx context.Context, client http.Client, req *http.Request) (int, []byte, error) {
	if ctx != nil {
		req = req.WithContext(ctx)
	}
	resp, err := client.Do(req)
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()

	if err != nil {
		return 0, nil, err
	}

	var body []byte
	done := make(chan struct{})
	go func() {
		body, err = ioutil.ReadAll(resp.Body)
		close(done)
	}()

	select {
	case <-ctx.Done():
		<-done
		err = resp.Body.Close()
		if err == nil {
			err = ctx.Err()
		}
	case <-done:
	}

	return resp.StatusCode, body, err
}
