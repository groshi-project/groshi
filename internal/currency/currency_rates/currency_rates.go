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

var errRateIsNotCached = errors.New("rate is not cached")
var errRateCacheIsExpired = errors.New("rate cache is expired")

// cacheIsExpired returns true if period longer than cacheTTL
// has passed since updatedAt.
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

// readRatesFromCache reads currency rates information from the database.
// Also checks if cache is expired and returns errCacheIsExpired in such case.
// Returns errRateIsNotCached if rates information could not be found.
func readRatesFromCache(baseCurrency string, v *database.CurrencyRates) error {
	if err := database.CurrencyRatesCol.FindOne(
		database.Context,
		bson.D{{"base_currency", baseCurrency}},
	).Decode(v); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) { // rates information was not found in the databse
			return errRateIsNotCached
		}
		return err // unexpected error
	}
	if cacheIsExpired(v.UpdatedAt) {
		return errRateCacheIsExpired // cache has already expired
	}
	return nil
}

// FetchRate is a complicated high-level function for
// fetching `baseCurrency` rate in `targetCurrency` units which:
//
// -> Tries to get rate from rates cache
//
//			<- If succeeded: returns rate you are looking for.
//			<- If cache is outdated (see cacheTTL): fetches rates using third-party API
//	    	and updates cache. Returns rate.
//			-> If rate information is not cached:
//					-> Tries to get rates from cache of reverse rate.
//							<- If succeeded: returns rate you are looking for.
//							-> If reverse rate is outdated: fetches reverse rate using
//							third-party API and updates its cache. Returns rate.
//							<- If reverse rate is not cached: fetches rate (not reverse rate)
//							Using third-party API, caches it and finally returns rate
func FetchRate(baseCurrency string, targetCurrency string) (float64, error) {
	ratesObject := database.CurrencyRates{}
	err := readRatesFromCache(baseCurrency, &ratesObject)
	if err == nil { // if rate is cached in the database and not expired
		rate, found := ratesObject.Rates[targetCurrency]
		if !found {
			return 0, ErrUnknownCurrency
		}
		return rate.(float64), nil

	} else if errors.Is(err, errRateCacheIsExpired) { // if rate cache is expired
		// fetch new rates, update cache and return fresh rate.
		rates, err := exchangerates.Client.GetRates(baseCurrency)
		if err != nil {
			return 0, err
		}

		if err := updateCachedRates(baseCurrency, rates); err != nil {
			return 0, err
		}

		rate, found := rates[targetCurrency]
		if !found {
			return 0, ErrUnknownCurrency
		}

		return rate.(float64), err

	} else if errors.Is(err, errRateIsNotCached) { // if rate is not cached in the database
		// then we will try to fetch rates for `targetCurrency` in order to get
		// rates for `baseCurrency` indirectly, using inverse rate

		inverseRatesObject := database.CurrencyRates{}
		// using `targetCurrency` (not `baseCurrency`) as base targetCurrency to try to fetch rates anyway indirectly

		err := readRatesFromCache(targetCurrency, &inverseRatesObject)
		if err == nil { // inverse rate is cached and is not expired
			// then simply calculate rate and return it!

			inverseRate, found := inverseRatesObject.Rates[baseCurrency]
			if !found {
				return 0, ErrUnknownCurrency
			}
			return 1 / inverseRate.(float64), nil

		} else if errors.Is(err, errRateCacheIsExpired) { // inverse rate is expired
			// then fetch new inverse rates, update cache and return fresh rate!

			inverseRates, err := exchangerates.Client.GetRates(targetCurrency)
			if err != nil {
				return 0, err
			}

			if err := updateCachedRates(targetCurrency, inverseRates); err != nil {
				return 0, err
			}

			inverseRate, found := inverseRates[targetCurrency]
			if !found {
				return 0, ErrUnknownCurrency
			}

			return 1 / inverseRate.(float64), err

		} else if errors.Is(err, errRateIsNotCached) {
			return 0, err

		} else { // inverse rate is not cached
			rates, err := exchangerates.Client.GetRates(baseCurrency)
			if err != nil {
				return 0, err
			}

			if err := updateCachedRates(baseCurrency, rates); err != nil {
				return 0, err
			}

			rate, found := rates[targetCurrency]
			if !found {
				return 0, ErrUnknownCurrency
			}

			return rate.(float64), err

		}
	} else { // other unexpected error
		return 0, err
	}
}
