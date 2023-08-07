package currency_rates

import (
	"errors"
	"github.com/jieggii/groshi/internal/currency/exchangerates"
	"github.com/jieggii/groshi/internal/database"
	"github.com/jieggii/groshi/internal/loggers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const cacheTTL = time.Hour * 24

var ErrUnknownCurrency = errors.New("unknown currency")

func cacheIsExpired(updatedAt time.Time) bool {
	sub := time.Now().Sub(updatedAt)
	if sub.Hours() >= cacheTTL.Hours() {
		return true
	}
	return false
}

// updateCachedRates updates cached rates.
func updateCachedRates(baseCurrency string, rates map[string]interface{}) error {
	_, err := database.CurrencyRatesCol.UpdateOne(
		database.Context,
		bson.D{{"base_currency", baseCurrency}},
		bson.D{
			{"rates", rates},
			{"updated_at", time.Now()},
		},
	)
	return err
}

// cacheRates saves rates to the database.
func cacheRates(baseCurrency string, rates map[string]interface{}) error {
	loggers.Info.Println("cacheRates was called")
	rateObject := database.CurrencyRates{
		ID:        primitive.NewObjectID(),
		Rates:     rates,
		UpdatedAt: time.Now(),
	}
	_, err := database.CurrencyRatesCol.InsertOne(database.Context, &rateObject)
	return err
}

// FetchRate TODO
func FetchRate(baseCurrency string, currency string) (float64, error) {
	currencyRateObject := database.CurrencyRates{}
	err := database.CurrencyRatesCol.FindOne(
		database.Context,
		bson.D{{"base_currency", baseCurrency}},
	).Decode(&currencyRateObject)
	if err == nil { // if rate is cached in the database
		loggers.Info.Println("using cached rates")
		if cacheIsExpired(currencyRateObject.UpdatedAt) { // if cache is outdated
			// fetch new rates, update cache and return fresh rate.
			rates, err := exchangerates.Client.GetRates(baseCurrency)
			if err != nil {
				return 0, err
			}
			if err := updateCachedRates(baseCurrency, rates); err != nil {
				return 0, err
			}
			rate, found := rates[currency]
			if !found {
				return 0, ErrUnknownCurrency
			}
			return rate.(float64), err
		} else { // cache is not outdated yet => return rate from cache
			rate, found := currencyRateObject.Rates[currency]
			if !found {
				return 0, ErrUnknownCurrency
			}
			return rate.(float64), err
		}

	} else if errors.Is(err, mongo.ErrNoDocuments) { // rate hasn't been cached in the database yet
		loggers.Info.Println("fetching rates")
		// fetch rates and cache them, return rate
		rates, err := exchangerates.Client.GetRates(baseCurrency)
		if err != nil {
			return 0, err
		}
		if err := cacheRates(baseCurrency, rates); err != nil {
			return 0, err
		}
		rate, found := rates[currency]
		if !found {
			return 0, ErrUnknownCurrency
		}
		return rate.(float64), nil

	} else { // other unexpected error
		return 0, err
	}
}
