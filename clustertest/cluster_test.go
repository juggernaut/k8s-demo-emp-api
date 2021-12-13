package clustertest

import (
	"encoding/json"
	"fmt"
	"k8s-demo-emp-api/api"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"testing"
)

func TestApiServer(t *testing.T) {

	apiBaseUrl := os.Getenv("API_BASE_URL")
	if apiBaseUrl == "" {
		t.Fatalf("API_BASE_URL must be set!")
	}
	client := &http.Client{}

	name := randName()
	age := randAge()
	var id string

	t.Run("test POST", func (t *testing.T) {
		resp, err := client.PostForm(fmt.Sprintf("%s/employees", apiBaseUrl), url.Values{
			"Name": {name},
			"Age": {strconv.Itoa(age)},
		})
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != 201 {
			t.Fatalf("Expected status code '201', got %d instead\n", resp.StatusCode)
		}

		decoder := json.NewDecoder(resp.Body)
		var emp api.Employee
		if err := decoder.Decode(&emp); err != nil {
			t.Fatal(err)
		}
		if emp.Name != name {
			t.Errorf("Unexpected name in response %s", emp.Name)
		}
		if emp.Age != age {
			t.Errorf("Unexpected age in response %d", emp.Age)
		}
		id = emp.Id
	})

	t.Run("test GET", func (t *testing.T) {
		resp, err := client.Get(fmt.Sprintf("%s/employees/%s", apiBaseUrl, id))
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != 200 {
			t.Fatalf("Expected status code '200', got %d instead\n", resp.StatusCode)
		}

		decoder := json.NewDecoder(resp.Body)
		var emp api.Employee
		if err := decoder.Decode(&emp); err != nil {
			t.Fatal(err)
		}
	})
}

func randName() string {
	var nb []byte
	for i := 0; i < 10; i++ {
		a := rand.Intn(26)
		nb = append(nb, 'a' + byte(a))
	}
	return string(nb)
}

func randAge() int {
	n := rand.Intn(50)
	return n + 20
}