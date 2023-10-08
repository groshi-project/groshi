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

// saveCache saves rates to the database.
func saveCache(v *database.CurrencyRates) error {
	_, err := database.CurrencyRatesCol.InsertOne(
		database.Context,
		v,
	)
	return err
}

// updateCache updates cached currency rates in the database.
func updateCache(v *database.CurrencyRates) error {
	_, err := database.CurrencyRatesCol.ReplaceOne(
		database.Context,
		bson.D{{"_id", v.ID}},
		v,
	)
	return err
}

// readFromCache returns rates stored in the database.
// If currency rates collection is empty, returns errRatesAreNotCached.
func readFromCache(v *database.CurrencyRates) error {
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

// readFromAPI fetches latest currency rates and updates fields of v with new data
// (also updates UpdatedAt field).
func readFromAPI(v *database.CurrencyRates) error {
	rates, err := exchangerates.Client.GetRates(baseCurrency)
	if err != nil {
		return err
	}
	v.Rates = rates

	return nil
}

// read reads up-to-date currency rates either from the database or third-part API,
// according to the cacheTTL. Also updates stored rates if needed.
func read() (map[string]interface{}, error) {
	rates := database.CurrencyRates{BaseCurrency: baseCurrency}
	err := readFromCache(&rates)

	if err != nil {
		if errors.Is(err, errRatesAreNotCached) { // if cache is not stored in the database
			// then read it using API and save
			loggers.Info.Println("fetching currency rates from third-party because rates are not cached")

			if err := readFromAPI(&rates); err != nil {
				return nil, err
			}

			// generate object id for the new document and set UpdatedAt field
			rates.ID = primitive.NewObjectID()
			rates.UpdatedAt = time.Now()

			if err := saveCache(&rates); err != nil {
				return nil, err
			}
		} else { // if unexpected error happened when fetching rates from the database
			return nil, err
		}

	} else { // if successfully fetched rates from the database
		// then check if cache is expired

		if time.Now().Sub(rates.UpdatedAt).Hours() > cacheTTL.Hours() { // if cache is expired
			// then read new rates using API and update stored cache

			loggers.Info.Println("fetching currency rates from third-party because cache is expired")
			if err := readFromAPI(&rates); err != nil {
				return nil, err
			}

			rates.UpdatedAt = time.Now()
			if err := updateCache(&rates); err != nil {
				return nil, err
			}
		}
	}

	// finally we have relevant rates! It's time to return them:
	return rates.Rates, nil
}

// getRate returns rate for given currency.
func getRate(currency string) (float64, error) {
	rates, err := read()
	if err != nil {
		return 0, err
	}

	rate, found := rates[currency]
	if !found {
		return 0, errors.New("unknown currency")
	}

	return rate.(float64), nil
}

// GetCurrencies returns slice of available ISO-4217 currency codes.
func GetCurrencies() ([]string, error) {
	currencies := make([]string, 0)

	rates, err := read()
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

	fromRate, err := getRate(fromCurrency)
	if err != nil {
		return 0, err
	}

	toRate, err := getRate(toCurrency)
	if err != nil {
		return 0, err
	}

	rate := toRate / fromRate

	return amount * rate, nil
}
