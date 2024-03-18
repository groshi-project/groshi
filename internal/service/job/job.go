package job

import (
	"github.com/groshi-project/groshi/internal/database"
)

// Job represents dependencies for the service jobs.
type Job struct {
	// database used to store and retrieve data.
	database *database.Database
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
	//		logger.Warning.Printf("could not parse rate '%s' of currency %s", currencyRate, currencyCode)
	//		continue
	//	}
	//	code := string(currencyCode)
	//
	//	currency := database.Currency{}
	//	err = s.Database.Client.NewSelect().Model(database.sampleCurrency).Where("code = ?", currencyCode).Scan(s.Database.Ctx, &currency)
	//	if err == nil {
	//		// update rate of existing currency:
	//		currency.Rate = rate
	//		if _, err := s.Database.Client.NewUpdate().Model(currency).WherePK().Exec(s.Database.Ctx); err != nil {
	//			logger.Warning.Println(err)
	//		}
	//	} else if errors.Is(err, sql.ErrNoRows) {
	//		// insert a new currency:
	//		currency.Code = code
	//		currency.Symbol = code // todo: get currency symbol (e.g. "$") from somewhere
	//		currency.Rate = rate
	//		if _, err := s.Database.Client.NewInsert().Model(currency).Exec(s.Database.Ctx); err != nil {
	//			logger.Warning.Println(err)
	//		}
	//	} else {
	//		logger.Warning.Println(err)
	//	}
	//}

	return nil
}
