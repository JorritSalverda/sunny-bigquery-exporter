package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {

	t.Run("ReturnsSuccess", func(t *testing.T) {

		if testing.Short() {
			t.Skip("skipping test in short mode.")
		}

		localPort := 8855
		host := "192.168.178.254"
		port := 9522
		user := "User"
		password := "0000"

		client, err := NewSunnyBoyClient(localPort, host, port, user, password)
		if !assert.Nil(t, err, "failed creating client") {
			return
		}

		// act
		success, err := client.Login()

		if !assert.Nil(t, err, "failed logging in") {
			return
		}

		assert.True(t, success)
	})
}
