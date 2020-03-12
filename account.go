package form3

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

const (
	organisationAccountsBasePath = "/v1/organisation/accounts"
)

type AccountRoot struct {
	Data Account `json:"data,omitempty"`
}

type Account struct {
	ID             string            `json:"id,omitempty"`
	OrganisationID string            `json:"organisation_id,omitempty"`
	CreatedOn      string            `json:"created_on,omitempty"`
	ModifiedOn     string            `json:"modified_on,omitempty"`
	Type           string            `json:"type,omitempty"`
	Version        int               `json:"version,omitempty"`
	Attributes     AccountAttributes `json:"attributes,omitempty"`
}

type AccountAttributes struct {
	AccountMatchingOptOut       bool     `json:"account_matching_opt_out,omitempty"`
	JointAccount                bool     `json:"joint_account,omitempty"`
	AccountClassification       string   `json:"account_classification,omitempty"`
	AccountNumber               string   `json:"account_number,omitempty"`
	AlternativeBankAccountNames []string `json:"alternative_bank_account_names,omitempty"`
	BankAccountName             string   `json:"bank_account_name,omitempty"`
	BankID                      string   `json:"bank_id,omitempty"`
	BankIDCode                  string   `json:"bank_id_code,omitempty"`
	BaseCurrency                string   `json:"base_currency,omitempty"`
	Bic                         string   `json:"bic,omitempty"`
	Country                     string   `json:"country,omitempty"`
	FirstName                   string   `json:"first_name,omitempty"`
	Iban                        string   `json:"iban,omitempty"`
	SecondaryIdentification     string   `json:"secondary_identification,omitempty"`
	Title                       string   `json:"title,omitempty"`
}

type AccountRequestRoot struct {
	Data AccountRequest `json:"data,omitempty"`
}

type AccountListRoot struct {
	Data  []Account `json:"data,omitempty"`
	Links struct {
		Next     string `json:"next,omitempty"`
		Previous string `json:"prev,omitempty"`
		First    string `json:"first,omitempty"`
		Last     string `json:"last,omitempty"`
		Self     string `json:"self,omitempty"`
	} `json:"links,omitempty"`
}

type AccountRequest struct {
	ID             string            `json:"id"`
	OrganisationID string            `json:"organisation_id"`
	Type           string            `json:"type,omitempty"`
	Attributes     AccountAttributes `json:"attributes"`
}

type AccountService struct {
	client *Client
}

func (a *AccountService) Create(createReq AccountRequest) (Account, error) {
	req, err := a.client.newRequest(http.MethodPost, &url.URL{Path: organisationAccountsBasePath}, AccountRequestRoot{createReq})
	if err != nil {
		return Account{}, err
	}

	result := AccountRoot{}
	err = a.client.do(req, &result)

	return result.Data, err
}

func (a *AccountService) Fetch(id string) (Account, error) {
	fetchAccountPath := fmt.Sprintf("%s/%s", organisationAccountsBasePath, id)
	req, err := a.client.newRequest(http.MethodGet, &url.URL{Path: fetchAccountPath}, nil)

	result := AccountRoot{}
	err = a.client.do(req, &result)

	return result.Data, err
}

func (a *AccountService) Delete(id string, version int) error {
	deleteAccountPath := fmt.Sprintf("%s/%s", organisationAccountsBasePath, id)

	deleteQuery := url.Values{
		"version": {strconv.Itoa(version)},
	}

	req, err := a.client.newRequest(http.MethodDelete, &url.URL{Path: deleteAccountPath, RawQuery: deleteQuery.Encode()}, nil)
	err = a.client.do(req, nil)

	return err
}

func (a *AccountService) List(options ListOptions) ([]Account, bool, error) {
	listQuery := a.getPagingQueryParams(options)

	req, err := a.client.newRequest(http.MethodGet, &url.URL{Path: organisationAccountsBasePath, RawQuery: listQuery.Encode()}, nil)

	result := AccountListRoot{}
	err = a.client.do(req, &result)

	return result.Data, result.Links.Next != "", err
}

func (a *AccountService) getPagingQueryParams(options ListOptions) url.Values {
	var (
		pageNumber string
		pageSize   string
	)

	if options.PageSize == 0 {
		pageSize = defaultPageSize
	} else {
		pageSize = strconv.Itoa(options.PageSize)
	}

	pageNumber = strconv.Itoa(options.Page)
	listQuery := url.Values{
		"page[number]": {pageNumber},
		"page[size]":   {pageSize},
	}
	return listQuery
}
