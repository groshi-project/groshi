package currency_rates

import (
	"errors"
	"github.com/jieggii/groshi/internal/currency/exchangerates"
	"github.com/jieggii/groshi/internal/database"
	"github.com/jieggii/groshi/internal/loggers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const baseCurrency string = "EUR"
const cacheTTL = time.Hour * 24

var errRatesAreNotCached = errors.New("rates are not cached")

// readRatesFromCache returns rates stored in the database.
// If currency rates collection is empty, returns errRatesAreNotCached.
func readRatesFromCache(v *database.CurrencyRates) error {
	cursor, err := database.CurrencyRatesCol.Find(
		database.Context,
		bson.M{},
	)
	if err != nil {
		return err
	}

	decoded := false
	for cursor.Next(database.Context) {
		if decoded {
			return errors.New("more than one documents in the currency rates collection")
		}
		if err := cursor.Decode(v); err != nil {
			return err
		}
		decoded = true
	}

	if !decoded {
		return errRatesAreNotCached
	}

	return nil
}

// readRatesFromAPI fetches latest currency rates and updates fields of v with new data
// (also updates UpdatedAt field).
func readRatesFromAPI(v *database.CurrencyRates) error {
	rates, err := exchangerates.Client.GetRates(baseCurrency)
	if err != nil {
		return err
	}

	v.Rates = rates
	v.UpdatedAt = time.Now()

	return nil
}

// cacheRates saves rates to the database.
func cacheRates(v *database.CurrencyRates) error {
	_, err := database.CurrencyRatesCol.InsertOne(
		database.Context,
		v,
	)
	return err
}

// updateCachedRates updates cached currency rates in the database.
func updateCachedRates(v *database.CurrencyRates) error {
	_, err := database.CurrencyRatesCol.ReplaceOne(
		database.Context,
		bson.D{{"id", v.ID}},
		v,
	)
	return err
}

// fetchRate fetches up-to-date currency rates either from the database or third-part API,
// according to the cacheTTL. Also updates cache if needed.
func fetchRate(currency string) (float64, error) {
	rates := database.CurrencyRates{BaseCurrency: baseCurrency}
	err := readRatesFromCache(&rates)

	if err != nil {
		if errors.Is(err, errRatesAreNotCached) { // if cache is not stored
			// then fetch it using API and store

			rates.ID = primitive.NewObjectID()
			if err := readRatesFromAPI(&rates); err != nil {
				return 0, err
			}
			if err := cacheRates(&rates); err != nil {
				return 0, err
			}
		} else { // if unexpected error happened when fetching rates from cache
			return 0, err
		}

	} else { // if successfully fetched rates from cache
		// then check if cache is expired

		if time.Now().Sub(rates.UpdatedAt).Hours() > cacheTTL.Hours() { // if cache is expired
			// then fetch new rates using API and update cache

			loggers.Info.Println("! cache is expired")
			rates.ID = primitive.NewObjectID()
			if err := readRatesFromAPI(&rates); err != nil {
				return 0, err
			}
			if err := updateCachedRates(&rates); err != nil {
				return 0, err
			}
		}
	}

	// finally we have relevant rates in `rates`! It's time to return the desired rate:
	rate, found := rates.Rates[currency]
	if !found {
		return 0, errors.New("unknown currency")
	}

	return rate.(float64), nil
}

func Convert(fromCurrency string, toCurrency string, amount float64) (float64, error) {
	if amount == 0 { // small optimization: if amount is 0, it will be 0 in any currency :)
		return 0, nil
	}

	fromRate, err := fetchRate(fromCurrency)
	if err != nil {
		return 0, err
	}

	toRate, err := fetchRate(toCurrency)
	if err != nil {
		return 0, err
	}

	rate := toRate / fromRate

	return amount * rate, nil
}
