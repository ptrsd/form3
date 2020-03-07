package form3

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

type AccountService struct {
	client *Client
}
