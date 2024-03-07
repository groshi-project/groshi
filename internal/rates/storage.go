package rates

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
)

var ErrRatesNotStored = errors.New("rates are not stored")

// ratesKey is name of they key which will contain rates information.
const ratesKey = "rates"

// Storage represents rates redis storage.
type Storage struct {
	redisCtx    context.Context
	redisClient *redis.Client
}

// NewStorage creates a new instance of [Storage] and returns pointer to it.
func NewStorage(redisClient *redis.Client) *Storage {
	return &Storage{
		redisClient: redisClient,
		redisCtx:    context.Background(),
	}
}

// GetRates returns rates stored in the redis storage.
// Returns [ErrRatesNotStored] as error if rates are not stored.
func (s *Storage) GetRates() (map[string]float64, error) {
	val, err := s.redisClient.Get(s.redisCtx, ratesKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) { // if key is expired or was not found
			return nil, ErrRatesNotStored
		}
		return nil, err
	}

	var rates map[string]float64
	if err := json.Unmarshal([]byte(val), &rates); err != nil {
		return nil, err
	}

	return rates, nil
}

// SaveRates saves rates to the redis storage.
func (s *Storage) SaveRates(rates map[string]float64) error {
	data, err := json.Marshal(rates)
	if err != nil {
		return err
	}

	if err := s.redisClient.Set(s.redisCtx, ratesKey, data, 0).Err(); err != nil {
		return err
	}
	return nil
}
