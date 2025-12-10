package remoteteecache

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/hyperledger/aries-framework-go/pkg/doc/did"
)

type RemoteTeeSecureCache struct {
	url string
}

func New(url string) (*RemoteTeeSecureCache, error) {
	cache_url := "http://155.54.205.92:5000"
	if url != "" {
		cache_url = url
	}
	return &RemoteTeeSecureCache{url: cache_url}, nil
}

func (c *RemoteTeeSecureCache) StoreDidDoc(document *did.Doc) error {
	if document == nil {
		return errors.New("value cannot be nil")
	}
	if document.ID == "" {
		return errors.New("didId cannot be empty")
	}

	serialDoc, err := document.MarshalJSON()
	if err != nil {
		fmt.Println("Error encoding JSON DIDDoc:", err)
		return err
	}

	data := map[string]string{
		"id":    "diddocument-" + document.ID,
		"value": base64.StdEncoding.EncodeToString(serialDoc),
	}

	serialData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error encoding JSON data:", err)
		return err
	}

	// Create a new POST request
	req, err := http.NewRequest("POST", c.url+"/store", bytes.NewBuffer(serialData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return err
	}

	// Print raw response
	//fmt.Println("Response status:", resp.Status)
	//fmt.Println("Response body:", string(body))

	// Optionally unmarshal JSON response
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err == nil {
		//fmt.Println("\nParsed response:")
		if result["status"] == "success" {
			return nil
		}

		//fmt.Println("Output:", result["output"])
		return errors.New("Status not success")
	}

	return err
}

func (c *RemoteTeeSecureCache) RetrieveDidDoc(didID string) (bool, *did.Doc, error) {
	if didID == "" {
		return false, nil, errors.New("didId cannot be empty")
	}

	// The user ID to load
	id := "diddocument-" + didID

	// Construct the URL with path parameter
	url := fmt.Sprintf("http://155.54.205.92:5000/load/%s", id)

	// Perform GET request
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error making GET request:", err)
		return false, nil, err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return false, nil, err
	}

	//fmt.Println("Response status:", resp.Status)
	//fmt.Println("Raw body:", string(body))

	// Define struct matching the expected JSON response
	type LoadResponse struct {
		ID    string `json:"id"`
		Value string `json:"value"`
	}

	var result LoadResponse
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return false, nil, err
	}

	//Value lleva NODATA -> no existía
	if result.Value == "NODATA" {
		return false, nil, nil
	}

	// Print parsed response
	//fmt.Println("\nParsed response:")
	//fmt.Println("ID:", result.ID)
	//fmt.Println("Value:", result.Value)
	data, err := base64.StdEncoding.DecodeString(result.Value)
	if err != nil {
		fmt.Println("Error decoding b64 json:", err)
		return false, nil, err
	}

	doc, err := did.ParseDocument(data)
	if err != nil {
		fmt.Println("Error parsing Doc:", err)
		return false, nil, err
	}
	return true, doc, nil
}

func (c *RemoteTeeSecureCache) RemoveAll() error {
	//NOOP as it is not supperted for now
	return nil
}
