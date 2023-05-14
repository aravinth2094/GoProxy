package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"

	"github.com/oov/socks5"
	"golang.org/x/crypto/bcrypt"
)

func Authenticate(username, password string) error {
	if err := authenticate(username, password); err != nil {
		return socks5.ErrAuthenticationFailed
	}
	return nil
}

func CheckAllowed(username, targetAddr string) error {
	return checkAllowed(username, targetAddr)
}

func authenticate(username, password string) error {
	item, err := fetchUser(username)
	if err != nil {
		return err
	}
	if err = bcrypt.CompareHashAndPassword([]byte(item["password"].(string)), []byte(password)); err != nil {
		return err
	}
	return nil
}

func fetchUser(username string) (map[string]interface{}, error) {
	// Read the JSON file
	filePath := "access.json"
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Define a slice to store the JSON array
	var jsonArray []map[string]interface{}

	// Unmarshal the JSON data into the slice
	err = json.Unmarshal(data, &jsonArray)
	if err != nil {
		return nil, err
	}

	// Process the JSON array
	for _, item := range jsonArray {
		// Access individual elements of the JSON array
		if item["username"] != username {
			continue
		}
		return item, nil
	}

	return nil, fmt.Errorf("user not found")
}

func checkAllowed(username, targetAddr string) error {
	user, err := fetchUser(username)
	if err != nil {
		return err
	}
	ip, port, err := net.SplitHostPort(targetAddr)
	if err != nil {
		return err
	}
	for _, access := range user["access"].([]interface{}) {
		portAllowed := false
		for _, allowedPort := range access.(map[string]interface{})["ports"].([]interface{}) {
			if port == allowedPort {
				portAllowed = true
			}
		}
		if !portAllowed {
			return fmt.Errorf("port not allowed")
		}
		network, err := fetchNetwork(access.(map[string]interface{})["network"].(string))
		if err != nil {
			return err
		}
		for _, node := range network["nodes"].([]interface{}) {
			if node.(map[string]interface{})["ip"].(string) == ip {
				for _, nodePort := range node.(map[string]interface{})["ports"].([]interface{}) {
					if port == nodePort {
						return nil
					}
				}
				return fmt.Errorf("port not listed")
			}
		}
		return fmt.Errorf("node not listed")
	}
	return fmt.Errorf("could not find appropriate network")
}

func fetchNetwork(network string) (map[string]interface{}, error) {
	// Read the JSON file
	filePath := "networks.json"
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Define a slice to store the JSON array
	var jsonArray []map[string]interface{}

	// Unmarshal the JSON data into the slice
	err = json.Unmarshal(data, &jsonArray)
	if err != nil {
		return nil, err
	}

	// Process the JSON array
	for _, item := range jsonArray {
		// Access individual elements of the JSON array
		if item["name"] != network {
			continue
		}
		return item, nil
	}

	return nil, fmt.Errorf("network not found")
}
