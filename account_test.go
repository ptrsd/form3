package form3

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"testing"
)

var baseURL string

func TestMain(m *testing.M) {
	baseURL = os.Getenv("APP_BASE_URL")
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	code := m.Run()
	os.Exit(code)
}

func TestAccountService_Create(t *testing.T) {
	client := NewClient(nil, baseURL)
	t.Run("When creating accounts with valid data then return new account", func(t *testing.T) {
		accountRequest, err := generateAccountWithAttributes(AccountAttributes{Country:"GB"})
		if err != nil {
			t.Errorf("error while generating minimal account, %s", err.Error())
			t.FailNow()
		}

		create, err := client.AccountService.Create(accountRequest)
		if err != nil {
			t.Errorf("create account returned with error %v", err.Error())
			t.FailNow()
		}

		equals := assertions{
			{actual: create.Type, expected: accountRequest.Type, name: "Type"},
			{actual: create.ID, expected: accountRequest.ID, name: "ID"},
			{actual: create.OrganisationID, expected: accountRequest.OrganisationID, name: "OrganisationID"},
			{actual: create.Attributes, expected: accountRequest.Attributes, name: "Attributes"},
		}
		thenEquals(t, equals)

		notEmpty := assertions{
			{actual: create.ModifiedOn, name: "ModifiedOn"},
			{actual: create.CreatedOn, name: "CreatedOn"},
		}
		thenNotEmpty(t, notEmpty)
	})

	t.Run("When creating duplicates then error", func(t *testing.T) {
		accountRequest, err := generateAccountWithAttributes(AccountAttributes{Country:"GB"})
		if err != nil {
			t.Errorf("error while generating minimal account, %s", err.Error())
			t.FailNow()
		}

		_, err = client.AccountService.Create(accountRequest)
		if err != nil {
			t.Errorf("create account returned with error %v", err.Error())
			t.FailNow()
		}

		_, err = client.AccountService.Create(accountRequest)
		if err == nil {
			t.Errorf("create should return error when creating duplicates")
			t.FailNow()
		}

		equals := assertions{
			{actual: err.Error(), expected: "Account cannot be created as it violates a duplicate constraint", name: "DuplicateError"},
		}
		thenEquals(t, equals)
	})
}

func TestAccountService_Fetch(t *testing.T) {
	client := NewClient(nil, baseURL)

	t.Run("When fetching existing account then return requested account", func(t *testing.T) {
		accountRequest, err := generateAccountWithAttributes(AccountAttributes{Country:"GB"})
		if err != nil {
			t.Errorf("error while generating minimal account, %s", err.Error())
			t.FailNow()
		}

		newAccount, err := client.AccountService.Create(accountRequest)
		if err != nil {
			t.Errorf("create account returned with error %v", err.Error())
			t.FailNow()
		}

		account, err := client.AccountService.Fetch(accountRequest.ID)
		if err != nil {
			t.Errorf("fetch account returned with error %v", err.Error())
			t.FailNow()
		}

		equals := assertions{
			{actual: account.Attributes, expected: accountRequest.Attributes, name: "FetchAccount.AccountAttributes"},
			{actual: account.ID, expected: accountRequest.ID, name: "FetchAccount.AccountID"},
			{actual: account.OrganisationID, expected: accountRequest.OrganisationID, name: "FetchAccount.OrganisationID"},
			{actual: account.CreatedOn, expected: newAccount.CreatedOn, name: "FetchAccount.CreatedOn"},
			{actual: account.ModifiedOn, expected: newAccount.ModifiedOn, name: "FetchAccount.ModifiedOn"},
		}
		thenEquals(t, equals)
	})

	t.Run("When fetching not existing account then fail", func(t *testing.T) {
		uuid, err := generateRandomUUID()
		if err != nil {
			t.Errorf("error while generating random uuid, %s", err.Error())
			t.FailNow()
		}

		_, err = client.AccountService.Fetch(uuid)
		if err == nil {
			t.Errorf("fetch account should returned with error")
			t.FailNow()
		}

		equals := assertions{
			{
				actual:   err.Error(),
				expected: fmt.Sprintf("record %s does not exist", uuid),
				name:     "NotExistingAccountErrorMessage",
			},
		}
		thenEquals(t, equals)
	})
}

func TestAccountService_Delete(t *testing.T) {
	client := NewClient(nil, baseURL)

	t.Run("When deleting existing account then success", func(t *testing.T) {
		account, err := givenMinimalAccount(client)
		if err != nil {
			t.Errorf("error while generating minimal account, %s", err.Error())
			t.FailNow()
		}

		err = client.AccountService.Delete(account.ID, 0)
		if err != nil {
			t.Errorf("error while deleting account %s", err.Error())
			t.FailNow()
		}

		ac, err := client.AccountService.Fetch(account.ID)
		if err == nil {
			t.Errorf("expecting error while fetching account")
			t.FailNow()
		}

		equals := assertions{
			{
				actual:   err.Error(),
				expected: fmt.Sprintf("record %s does not exist", account.ID),
				name:     "NotExistingAccountErrorMessage",
			},
			{
				actual:   ac,
				expected: Account{},
				name:     "Account should be empty",
			},
		}
		thenEquals(t, equals)
	})
}

func TestAccountService_List(t *testing.T) {
	client := NewClient(nil, baseURL)
	clean(t, client)

	t.Run("Given no accounts", func(t *testing.T) {
		t.Run("When no accounts then list is empty", func(t *testing.T) {
			list, hasNext := whenListingAccountsWith(t, client, ListOptions{PageSize: 5})
			thenEquals(t, assertions{
				{actual: len(list), expected: 0, name: "ListLength"},
				{actual: hasNext, expected: false, name: "hasNext"},
			})
		})
	})

	t.Run("Given 101 accounts", func(t *testing.T) {
		givenAccounts(t, client, 101)

		t.Run("When using default pagination only 100 elements returned from 101 elements", func(t *testing.T) {
			list, hasNext := whenListingAccountsWith(t, client, ListOptions{})
			thenEquals(t, assertions{
				{actual: len(list), expected: 100, name: "ListLength"},
				{actual: hasNext, expected: true, name: "HasNext"},
			})
		})
		t.Run("When using default page size but asking for second page only 1 element returned from 101 elements", func(t *testing.T) {
			list, hasNext := whenListingAccountsWith(t, client, ListOptions{Page: 1})
			thenEquals(t, assertions{
				{actual: len(list), expected: 1, name: "ListLength"},
				{actual: hasNext, expected: false, name: "HasNext"},
			})
		})
		t.Run("When requesting first page with 5 elements only 5 elements returned from 101 elements", func(t *testing.T) {
			list, hasNext := whenListingAccountsWith(t, client, ListOptions{PageSize: 5})
			thenEquals(t, assertions{
				{actual: len(list), expected: 5, name: "ListLength"},
				{actual: hasNext, expected: true, name: "HasNext"},
			})
		})
	})
}

func whenListingAccountsWith(t *testing.T, client *Client, opts ListOptions) ([]Account, bool) {
	list, hasNext, err := client.AccountService.List(opts)
	if err != nil {
		t.Errorf("error while listing accounts, %s", err.Error())
		t.FailNow()
	}

	return list, hasNext
}

func givenAccounts(t *testing.T, client *Client, numberOfAccounts int) {
	for idx := 0; idx < numberOfAccounts; idx++ {
		_, err := givenMinimalAccount(client)
		if err != nil {
			t.Errorf("error while generating minimal account, %s", err.Error())
			t.FailNow()
		}
	}
}

func clean(t *testing.T, client *Client) {
	hasNext := true

	for hasNext {
		var list []Account
		var err error

		list, hasNext, err = client.AccountService.List(ListOptions{})
		if err != nil {
			t.Errorf("error while listing accounts, %s", err.Error())
			t.FailNow()
		}

		for _, acc := range list {
			if err := client.AccountService.Delete(acc.ID, acc.Version); err != nil {
				t.Errorf("error while deleting account, %s", err.Error())
				t.FailNow()
			}
		}
	}
}

func givenMinimalAccount(client *Client) (Account, error) {
	req, err := generateAccountWithAttributes(AccountAttributes{Country: "GB"})
	if err != nil {
		return Account{}, err
	}

	account, err := client.AccountService.Create(req)
	if err != nil {
		return Account{}, err
	}

	return account, nil
}

func generateAccountWithAttributes(attrs AccountAttributes) (AccountRequest, error) {
	accountRequest, err := generateAccountRequest()
	if err != nil {
		return AccountRequest{}, err
	}

	accountRequest.Attributes = attrs
	return accountRequest, nil
}

func generateAccountRequest() (AccountRequest, error) {
	id, orgID, err := generateIDs()
	if err != nil {
		return AccountRequest{}, err
	}

	accountRequest := AccountRequest{
		ID:             id,
		OrganisationID: orgID,
		Type:           "accounts",
	}

	return accountRequest, nil
}

func generateIDs() (id, orgID string, err error) {
	if id, err = generateRandomUUID(); err == nil {
		if orgID, err = generateRandomUUID(); err == nil {
			return id, orgID, err
		}
	}

	return "", "", err
}

func generateRandomUUID() (string, error) {
	elements, err := generateUUIDElements(8)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v%v-%v-%v-%v-%v%v%v", elements...), nil
}

func generateUUIDElements(numberOfElements int) ([]interface{}, error) {
	uuidElements := make([]interface{}, numberOfElements)
	for idx := 0; idx < numberOfElements; idx++ {

		uuidElement, err := generateHex(2)
		if err != nil {
			return nil, err
		}

		uuidElements[idx] = uuidElement
	}

	return uuidElements, nil
}

func generateHex(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}
