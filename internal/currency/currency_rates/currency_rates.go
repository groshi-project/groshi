package currency_rates

import (
	"errors"
	"github.com/groshi-project/groshi/internal/currency/exchangerates"
	"github.com/groshi-project/groshi/internal/database"
	"github.com/groshi-project/groshi/internal/loggers"
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
func updateRatesFromAPI(v *database.CurrencyRates) error {
	rates, err := exchangerates.Client.GetRates(baseCurrency)
	if err != nil {
		return err
	}
	v.Rates = rates

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
func replaceCachedRates(v *database.CurrencyRates) error {
	_, err := database.CurrencyRatesCol.ReplaceOne(
		database.Context,
		bson.D{{"_id", v.ID}},
		v,
	)
	return err
}

// fetchRates fetches up-to-date currency rates either from the database or third-part API,
// according to the cacheTTL. Also updates cache if needed.
func fetchRates() (map[string]interface{}, error) {
	ratesDocument := database.CurrencyRates{BaseCurrency: baseCurrency}
	err := readRatesFromCache(&ratesDocument)

	if err != nil {
		if errors.Is(err, errRatesAreNotCached) { // if cache is not stored
			// then fetch it using API and store
			loggers.Info.Println("fetching currency rates from third-party because rates are not cached")

			if err := updateRatesFromAPI(&ratesDocument); err != nil {
				return nil, err
			}

			// generate object id for the new document and set UpdatedAt field
			ratesDocument.ID = primitive.NewObjectID()
			ratesDocument.UpdatedAt = time.Now()

			if err := cacheRates(&ratesDocument); err != nil {
				return nil, err
			}
		} else { // if unexpected error happened when fetching rates from cache
			return nil, err
		}

	} else { // if successfully fetched rates from cache
		// then check if cache is expired

		if time.Now().Sub(ratesDocument.UpdatedAt).Hours() > cacheTTL.Hours() { // if cache is expired
			// then fetch new rates using API and update cache

			loggers.Info.Println("fetching currency rates from third-party because cache is expired")
			if err := updateRatesFromAPI(&ratesDocument); err != nil {
				return nil, err
			}

			ratesDocument.UpdatedAt = time.Now()

			if err := replaceCachedRates(&ratesDocument); err != nil {
				return nil, err
			}
		}
	}

	// finally we have relevant rates! It's time to return them:
	return ratesDocument.Rates, nil
}

// fetchRate returns rate for desired currency.
func fetchRate(currency string) (float64, error) {
	rates, err := fetchRates()
	if err != nil {
		return 0, err
	}

	rate, found := rates[currency]
	if !found {
		return 0, errors.New("unknown currency")
	}

	return rate.(float64), nil
}

// FetchCurrencies returns slice of available ISO-4217 currency codes.
func FetchCurrencies() ([]string, error) {
	currencies := make([]string, 0)

	rates, err := fetchRates()
	if err != nil {
		return currencies, err
	}
	for currency, _ := range rates {
		currencies = append(currencies, currency)
	}
	return currencies, nil
}

// Convert converts amount from one currency units to another.
// ?todo: make `amount` a param of type int.
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
