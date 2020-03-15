# FORM3-Client

_Author: Piotr Siudy_

## Usage

### Install

For latest version:

```shell script
go get github.com/ptrsd/form3
```

### Example

#### Create account

This example creates a new account, and prints the ID of the created account.

```go
package main

import (
	"fmt"
	"github.com/ptrsd/form3"
)

func main() {
	client := form3.NewDefaultClient(nil)
	account, _ := client.AccountService.Create(form3.AccountRequest{
		ID: "5b438472-e8f7-4ce5-a189-2968e6f8f62e",
		OrganisationID:"bfb86474-e82a-497f-8b65-c8a2c7f2fa44",
		Attributes:form3.AccountAttributes{
			Country:"GB",
		}})

	fmt.Printf("%v", account.ID)
}
```

Output:

```go
5b438472-e8f7-4ce5-a189-2968e6f8f62e
```

#### Pagination

This example shows how to list accounts, and how to use pagination for a list operation.

```go
package main

import (
	"fmt"
	"github.com/ptrsd/form3"
)

func main() {
	var (
		page     int
		accounts []form3.Account
		err      error

		hasNext = true
		client  = form3.NewDefaultClient(nil)
	)

	for hasNext {
		accounts, hasNext, err = client.AccountService.List(form3.ListOptions{
			PageSize: 2,
			Page:     page,
		})

		if err != nil {
			fmt.Errorf(err.Error())
		}

		fmt.Println(fmt.Sprintf("page no. %d", page))
		for _, account := range accounts {
			fmt.Println(account.ID)
		}

		page++
	}
}
```

Output:

```go
page no. 0
3ed7b128-b770-22da-cc61-d3f75c51a881
1a7ba8e0-bbb5-e4ed-44c5-d6a22192b148
page no. 1
f30fb8b6-52a7-3fe7-816e-503e0dee3890
928b105f-0277-cc0e-b194-9135345fcb43
page no. 2
...
```

## For developers

### Prerequisites

* Go v. 1.14
* Docker
* docker-compose

### How to

Makefile consist of following stages:

* docs - generates pdf from this file. The output file can be found in ```./build/output/README.pdf```.
* lint - runs golangci-lint against the code.
* test - runs tests against test environment set up by docker-compose.

You can also run tests by using ```docker-compose up``` in the root directory of the project.