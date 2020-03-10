package form3

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"testing"
)

func TestAccountService_whenCreateRequestIsValidThenReturnNewAccount(t *testing.T) {
	client := NewDefaultClient(nil)
	id, orgID, err := generateIDs()

	if err != nil {
		t.Errorf("error while generating random UUID, %s", err.Error())
	}

	accountRequest := AccountRequest{
		ID:             id,
		OrganisationID: orgID,
		Type:           "accounts",
		Attributes: AccountAttributes{
			Country: "GB",
		},
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
	id, orgID, err := generateIDs()

	if err != nil {
		t.Errorf("error while generating random UUID, %s", err.Error())
	}

	accountRequest := AccountRequest{
		ID:             id,
		OrganisationID: orgID,
		Type:           "accounts",
		Attributes: AccountAttributes{
			Country: "GB",
		},
	}

	_, err = client.AccountService.Create(accountRequest)
	if err != nil {
		t.Errorf("create account returned with error %v", err.Error())
	}

	_, err = client.AccountService.Create(accountRequest)
	if err == nil {
		t.Errorf("create should return error when creating duplicates")
	}
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
