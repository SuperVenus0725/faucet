package cosmosfaucet

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"faucet/cmd/config"

	"github.com/ignite/cli/ignite/pkg/xhttp"
)

type TransferRequest struct {
	// AccountAddress to request for coins.
	AccountAddress string `json:"address"`

	// reCaptcha response.
	ReCaptchaResponse string `json:"response"`

	// Coins that are requested.
	// default ones used when this one isn't provided.
	Coins []string `json:"coins"`
}

func NewTransferRequest(accountAddress string, coins []string) TransferRequest {
	return TransferRequest{
		AccountAddress: accountAddress,
		Coins:          coins,
	}
}

type TransferResponse struct {
	Error string `json:"error,omitempty"`
}

func (f Faucet) faucetHandler(w http.ResponseWriter, r *http.Request) {
	var req TransferRequest
	cookie_captcha, err := r.Cookie("response")
	if err != nil {
		responseError(w, http.StatusBadRequest, err)
		return
	}

	if err = cookie_captcha.Valid(); err != nil {
		responseError(w, http.StatusBadRequest, err)
		return
	}

	fmt.Println("00000Passed here00000")
	fmt.Println(cookie_captcha.Value)
	result, err := f.validateReCAPTCHA(cookie_captcha.Value)

	fmt.Println("Verify result")
	fmt.Println(result)

	if !result || err != nil {
		responseError(w, http.StatusBadRequest, err)
		return
	}

	// decode request into req.
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseError(w, http.StatusBadRequest, err)
		return
	}

	// determine coins to transfer.
	coins, err := f.coinsFromRequest(req)
	if err != nil {
		responseError(w, http.StatusBadRequest, err)
		return
	}

	// try performing the transfer
	if errCode, err := f.Transfer(r.Context(), req.AccountAddress, coins); err != nil {
		if err == context.Canceled {
			return
		}
		responseError(w, (int)(errCode), err)
	} else {
		responseSuccess(w)
	}
}

func (f Faucet) validateReCAPTCHA(recaptchaResponse string) (bool, error) {
	// Check this URL verification details from Google
	// https://developers.google.com/recaptcha/docs/verify
	req, err := http.PostForm(f.configuration.ReCAPTCHA_VerifyURL, url.Values{
		"secret":   {f.configuration.ReCAPTCHA_ServerKey},
		"response": {recaptchaResponse},
	})
	if err != nil { // Handle error from HTTP POST to Google reCAPTCHA verify server
		return false, err
	}
	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body) // Read the response from Google
	if err != nil {
		return false, err
	}

	var googleResponse config.GoogleRecaptchaResponse
	err = json.Unmarshal(body, &googleResponse) // Parse the JSON response from Google
	if err != nil {
		return false, err
	}
	return true, nil
}

// FaucetInfoResponse is the faucet info payload.
type FaucetInfoResponse struct {
	// IsAFaucet indicates that this is a faucet endpoint.
	// useful for auto discoveries.
	IsAFaucet bool `json:"is_a_faucet"`

	// ChainID is chain id of the chain that faucet is running for.
	ChainID string `json:"chain_id"`
}

func (f Faucet) faucetInfoHandler(w http.ResponseWriter, r *http.Request) {
	xhttp.ResponseJSON(w, http.StatusOK, FaucetInfoResponse{
		IsAFaucet: true,
		ChainID:   f.chainID,
	})
}

// coinsFromRequest determines tokens to transfer from transfer request.
func (f Faucet) coinsFromRequest(req TransferRequest) (sdk.Coins, error) {
	if len(req.Coins) == 0 {
		return f.coins, nil
	}

	var coins []sdk.Coin
	for _, c := range req.Coins {
		coin, err := sdk.ParseCoinNormalized(c)
		if err != nil {
			return nil, err
		}
		coins = append(coins, coin)
	}

	return coins, nil
}

func responseSuccess(w http.ResponseWriter) {
	xhttp.ResponseJSON(w, http.StatusOK, TransferResponse{})
}

func responseError(w http.ResponseWriter, code int, err error) {
	xhttp.ResponseJSON(w, code, TransferResponse{
		Error: err.Error(),
	})
}
