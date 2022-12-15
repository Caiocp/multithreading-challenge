package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type ApiCep struct {
	Code       string `json:"code"`
	State      string `json:"state"`
	City       string `json:"city"`
	District   string `json:"district"`
	Address    string `json:"address"`
	Status     int    `json:"status"`
	Ok         bool   `json:"ok"`
	StatusText string `json:"statusText"`
}

type ViaCep struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func main() {
	cep := os.Args[1]

	if len(cep) != 8 {
		log.Fatal("Invalid CEP format! Must be 8 digits without hyphen.")
	}

	viaCepChannel := make(chan ViaCep)
	apiCepChannel := make(chan ApiCep)

	formattedCep := cep[0:5] + "-" + cep[5:8]

	go func() {
		// time.Sleep(time.Second * 2)
		viaCepChannel <- getViaCepData(cep)
	}()

	go func() {
		// time.Sleep(time.Second * 2)
		apiCepChannel <- getApiCepData(formattedCep)
	}()

	select {
	case response := <-viaCepChannel:
		fmt.Println("ViaCep |", response)

	case response := <-apiCepChannel:
		fmt.Println("ApiCep |", response)

	case <-time.After(time.Second * 1):
		log.Fatal("Timeout")
	}
}

func getViaCepData(cep string) ViaCep {
	req, err := http.NewRequest("GET", "https://viacep.com.br/ws/"+cep+"/json/", nil)
	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var viaCep ViaCep
	json.Unmarshal(body, &viaCep)

	return viaCep
}

func getApiCepData(cep string) ApiCep {
	req, err := http.NewRequest("GET", "https://cdn.apicep.com/file/apicep/"+cep+".json", nil)
	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var apiCep ApiCep
	json.Unmarshal(body, &apiCep)

	return apiCep
}
