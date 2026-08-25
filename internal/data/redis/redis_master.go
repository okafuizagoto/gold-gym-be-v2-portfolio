package goldgym

import (
	"context"
	"encoding/json"
	"time"

	"github.com/pkg/errors"
)

// JSON BASED
func (d Data) AddToRedis(ctx context.Context, data interface{}, key string, time time.Duration) (err error) {
	jsoned, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "[addToRedis]")
	}

	// return d.rdb.Set(ctx, key, jsoned, 3600*time.Second).Err()
	return d.rdb.Set(ctx, key, jsoned, time).Err()
}

// JSON BASED
func (d Data) GetFromRedis(ctx context.Context, key string, dest interface{}) (err error) {
	result, err := d.rdb.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}

	return json.Unmarshal(result, dest)
}

func (d Data) DeleteFromRedis(ctx context.Context, key string) error {
	result, err := d.rdb.Del(ctx, key).Result()
	if err != nil {
		return errors.Wrap(err, "[DeleteFromRedis]")
	}

	// result == 0 artinya key tidak ditemukan (token sudah expired atau tidak valid)
	if result == 0 {
		return errors.New("token not found or already expired")
	}

	return nil
}
