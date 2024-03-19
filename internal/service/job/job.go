package job

import (
	"github.com/groshi-project/groshi/internal/database"
	"log"
)

// Job represents dependencies for the service jobs.
type Job struct {
	// database used to store and retrieve data.
	database *database.Database

	// errLogger used to log warnings and errors.
	errLogger *log.Logger
}

// New creates a new instance of [Job] and returns pointer to it.
func New(database *database.Database) *Job {
	return &Job{database: database}
}

// UpdateCurrencies updates currencies and their rates
// using information from https://ecb.europa.eu.
func (j *Job) UpdateCurrencies() error {
	//ecbRates, err := ecbrates.New()
	//if err != nil {
	//	return err
	//}
	//
	//for currencyCode, currencyRate := range ecbRates.Rate {
	//	rate, err := strconv.ParseFloat(currencyRate.(string), 64)
	//	if err != nil {
	//		j.errLogger.Printf("could not parse rate '%s' of currency %s", currencyRate, currencyCode)
	//		continue
	//	}
	//	code := string(currencyCode)
	//
	//	currency := &database.Currency{}
	//	if err := j.database.SelectCurrencyByCode(code, currency); err != nil {
	//		if errors.Is(err, sql.ErrNoRows) {
	//			currency.Code = code
	//			currency.Symbol = code // todo: get currency symbol (e.g. "$") from somewhere
	//			currency.Rate = rate
	//			if _, err := s.Database.Client.NewInsert().Model(currency).Exec(s.Database.Ctx); err != nil {
	//				j.errLogger.Printf("could not create new currency %s: %s", code, err)
	//			}
	//			continue
	//		}
	//		j.errLogger.Printf("could not fetch currency %s: %s", code, err)
	//		continue
	//	}
	//
	//	currency.Rate = rate
	//	if _, err := s.Database.Client.NewUpdate().Model(currency).WherePK().Exec(s.Database.Ctx); err != nil {
	//		j.errLogger.Printf("could not update currency %s: %s", code, err)
	//	}
	//
	//}
	//
	return nil
}
