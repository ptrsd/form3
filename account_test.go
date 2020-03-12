package form3

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"testing"
)

func TestAccountService_whenCreateRequestIsValidThenReturnNewAccount(t *testing.T) {
	client := NewDefaultClient(nil)
	accountRequest, err := generateMinimalAccount()

	if err != nil {
		t.Errorf("error while generating minimal account, %s", err.Error())
	}

	create, err := client.AccountService.Create(accountRequest)
	if err != nil {
		t.Errorf("create account returned with error %v", err.Error())
	}

	equals := assertions{
		{actual: create.Type, expected: accountRequest.Type, name: "Type"},
		{actual: create.ID, expected: accountRequest.ID, name: "ID"},
		{actual: create.OrganisationID, expected: accountRequest.OrganisationID, name: "OrganisationID"},
		{actual: create.Attributes, expected: accountRequest.Attributes, name: "Attributes"},
	}
	assertEquals(t, equals)

	notEmpty := assertions{
		{actual: create.ModifiedOn, name: "ModifiedOn"},
		{actual: create.CreatedOn, name: "CreatedOn"},
	}
	assertNotEmpty(t, notEmpty)
}

func TestAccountService_whenCreatingDuplicatesThenError(t *testing.T) {
	client := NewDefaultClient(nil)
	accountRequest, err := generateMinimalAccount()
	if err != nil {
		t.Errorf("error while generating minimal account, %s", err.Error())
	}

	_, err = client.AccountService.Create(accountRequest)
	if err != nil {
		t.Errorf("create account returned with error %v", err.Error())
	}

	_, err = client.AccountService.Create(accountRequest)
	if err == nil {
		t.Errorf("create should return error when creating duplicates")
	}

	equals := assertions{
		{actual: err.Error(), expected: "Account cannot be created as it violates a duplicate constraint", name: "CreateAccount.DuplicateError"},
	}
	assertEquals(t, equals)
}

func TestAccountService_whenFetchingExistingAccountThenSuccess(t *testing.T) {
	client := NewDefaultClient(nil)
	accountRequest, err := generateMinimalAccount()
	if err != nil {
		t.Errorf("error while generating minimal account, %s", err.Error())
	}

	newAccount, err := client.AccountService.Create(accountRequest)
	if err != nil {
		t.Errorf("create account returned with error %v", err.Error())
	}

	account, err := client.AccountService.Fetch(accountRequest.ID)
	if err != nil {
		t.Errorf("fetch account returned with error %v", err.Error())
	}

	equals := assertions{
		{actual: account.Attributes, expected: accountRequest.Attributes, name: "FetchAccount.AccountAttributes"},
		{actual: account.ID, expected: accountRequest.ID, name: "FetchAccount.AccountID"},
		{actual: account.OrganisationID, expected: accountRequest.OrganisationID, name: "FetchAccount.OrganisationID"},
		{actual: account.CreatedOn, expected: newAccount.CreatedOn, name: "FetchAccount.CreatedOn"},
		{actual: account.ModifiedOn, expected: newAccount.ModifiedOn, name: "FetchAccount.ModifiedOn"},
	}

	assertEquals(t, equals)
}

func TestAccountService_whenFetchingNotExistingAccountThenFail(t *testing.T) {
	client := NewDefaultClient(nil)
	accountRequest, err := generateMinimalAccount()
	if err != nil {
		t.Errorf("error while generating minimal account, %s", err.Error())
	}

	_, err = client.AccountService.Create(accountRequest)
	if err != nil {
		t.Errorf("create account returned with error %v", err.Error())
	}

	uuid, err := generateRandomUUID()
	if err != nil {
		t.Errorf("error while generating random uuid, %s", err.Error())
	}

	_, err = client.AccountService.Fetch(uuid)
	if err == nil {
		t.Errorf("fetch account should returned with error")
	}

	equals := assertions{
		{
			actual:   err.Error(),
			expected: fmt.Sprintf("record %s does not exist", uuid),
			name:     "FetchAccount.NotExistingAccountErrorMessage",
		},
	}
	assertEquals(t, equals)
}

func TestAccountService_Delete(t *testing.T) {
	client := NewDefaultClient(nil)
	accountRequest, err := generateMinimalAccount()
	if err != nil {
		t.Errorf("error while generating minimal account, %s", err.Error())
	}

	_, err = client.AccountService.Create(accountRequest)
	if err != nil {
		t.Errorf("create account returned with error %v", err.Error())
	}

	err = client.AccountService.Delete(accountRequest.ID, 0)
	if err != nil {
		t.Errorf(err.Error())
	}

	ac, err := client.AccountService.Fetch(accountRequest.ID)
	if err == nil {
		t.Errorf("%#v", ac)
	}

	equals := assertions{
		{
			actual:   err.Error(),
			expected: fmt.Sprintf("record %s does not exist", accountRequest.ID),
			name:     "DeleteAccount.NotExistingAccountErrorMessage",
		},
	}
	assertEquals(t, equals)
}

func generateMinimalAccount() (AccountRequest, error) {
	return generateAccountWithAttributes(AccountAttributes{Country: "GB"})
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
