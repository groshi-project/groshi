package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/groshi-project/groshi/internal/currency/currency_rates"
	"github.com/groshi-project/groshi/internal/handlers/util"
)

// CurrenciesRead returns slice of available currency codes in ISO-4217 format.
//
//	@summary		get array of available currency codes
//	@description	Returns array of available currency codes in ISO-4217 format.
//	@tags			currencies
//	@accept			json
//	@produce		json
//	@success		200	{object}	[]string	"Array of currency codes in ISO-4217 format is returned."
//	@router			/currencies [get]
func CurrenciesRead(c *gin.Context) {
	// IMPORTANT NOTE:
	// currency_rates.FetchCurrencies returns an array of supported currencies using
	// either cache or third-party (if cache is expired).
	//
	// It is also important to mention, that currency_rates.FetchCurrencies is also used once
	// in the beginning of the runtime to initialize `currency` and `optional_currency` validators.
	//
	// Based on the two previous facts the following situation is possible (yet highly improbable):
	// 1. `currency` and `optional_currency` validators are initialized in the beginning of the runtime.
	// 2. The third party updates list of supported currencies (deletes existing or adds new currencies).
	// 2. currency_rates.cacheTTL time passes and currencies cache is expired.
	// 3. ReadCurrencies groshi API method is triggered, new list of supported currencies is fetched
	//	  from the third party and returned to an API user.
	// 4. An API user uses currency, that was present in the list of supported currencies returned to him,
	//    but is not known by validators.
	// 5. `currency` or `optional_currency` validator fails, because slice of known currencies has not been updated.
	//
	// BUT: For now I see no point to fix that because the third party has a stable list of supported currencies,
	//      and it will unlikely be changed.
	currencies, err := currency_rates.FetchCurrencies()
	if err != nil {
		util.AbortWithStatusInternalServerError(c, err)
	}
	util.ReturnSuccessfulResponse(c, currencies)
}
