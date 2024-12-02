package ssh

import (
	"fmt"
	"strings"
	"time"

	gossh "golang.org/x/crypto/ssh"
)

type ConnInfo struct {
	User        string        `json:"user"`
	Addr        string        `json:"addr"`
	Port        int           `json:"port"`
	AuthMode    string        `json:"authMode"`
	Password    string        `json:"password"`
	PrivateKey  []byte        `json:"privateKey"`
	PassPhrase  []byte        `json:"passPhrase"`
	DialTimeOut time.Duration `json:"dialTimeOut"`

	Client     *gossh.Client  `json:"client"`
	Session    *gossh.Session `json:"session"`
	LastResult string         `json:"lastResult"`
}

// NewClient is a method of the ConnInfo struct that creates and configures a new SSH client connection.
// It sets up the necessary configuration for the SSH client based on the information provided in the ConnInfo struct.
// If the address contains a colon, it encloses it in square brackets to handle IPv6 addresses properly.
// Then it configures the authentication method, timeout, and other settings before attempting to establish the SSH connection.
func (c *ConnInfo) NewClient() (*ConnInfo, error) {
	// If the address in ConnInfo contains a colon, format it by enclosing it in square brackets.
	// This is a common way to handle IPv6 addresses to ensure proper parsing and connection.
	if strings.Contains(c.Addr, ":") {
		c.Addr = fmt.Sprintf("[%s]", c.Addr)
	}

	// Create a new instance of gossh.ClientConfig to hold the configuration settings for the SSH client.
	// This struct will be populated with various settings like authentication, timeout, etc.
	config := &gossh.ClientConfig{}

	// Set default values for the gossh.ClientConfig. These defaults might include things like default algorithms
	// for encryption, key exchange, etc. depending on the gossh library implementation.
	config.SetDefaults()

	// Construct the full address string in the format "address:port" using the address and port from the ConnInfo struct.
	// This will be used to specify the destination for the SSH connection.
	addr := fmt.Sprintf("%s:%d", c.Addr, c.Port)

	// Set the username for the SSH connection. The username is retrieved from the ConnInfo struct.
	config.User = c.User

	// Check the authentication mode specified in the ConnInfo struct.
	if c.AuthMode == "password" {
		// If the authentication mode is "password", set up the authentication method to use a password.
		// Create a gossh.AuthMethod using the password from the ConnInfo struct.
		config.Auth = []gossh.AuthMethod{gossh.Password(c.Password)}
	} else {
		// If the authentication mode is not "password", assume it's using a private key for authentication.
		// First, create a signer using the private key and passphrase (if any) from the ConnInfo struct.
		signer, err := makePrivateKeySigner(c.PrivateKey, c.PassPhrase)
		if err != nil {
			// If there is an error creating the signer, return nil and the error.
			// This will propagate the error back to the caller of this function.
			return nil, err
		}
		// Set up the authentication method to use the public key corresponding to the created signer.
		config.Auth = []gossh.AuthMethod{gossh.PublicKeys(signer)}
	}

	// Check if the dial timeout value in the ConnInfo struct is zero.
	if c.DialTimeOut == 0 {
		// If it's zero, set a default dial timeout of 5 seconds.
		// This determines how long the function will wait when attempting to establish the SSH connection.
		c.DialTimeOut = 5 * time.Second
	}

	// Set the timeout value for the SSH client connection. The timeout value is retrieved from the ConnInfo struct.
	config.Timeout = c.DialTimeOut

	// Set the host key callback function. In this case, it's set to insecurely ignore the host key.
	// This is not recommended for production use as it bypasses host key verification, but might be useful for testing or
	// in some specific scenarios where the host key is known to be valid in another way.
	config.HostKeyCallback = gossh.InsecureIgnoreHostKey()

	// Determine the protocol to use for the SSH connection. If the address contains a colon, assume it's an IPv6 address
	// and set the protocol to "tcp6". Otherwise, use the default "tcp" protocol.
	proto := "tcp"
	if strings.Contains(c.Addr, ":") {
		proto = "tcp6"
	}

	// Attempt to establish the SSH client connection using the configured protocol, address, and client configuration.
	client, err := gossh.Dial(proto, addr, config)
	if nil != err {
		// If there is an error establishing the connection, return the current ConnInfo struct and the error.
		// This allows the caller to handle the error appropriately and potentially retry the connection.
		return c, err
	}

	// If the connection is successfully established, assign the created SSH client to the ConnInfo struct's Client field.
	c.Client = client

	// Return the updated ConnInfo struct and a nil error to indicate success.
	return c, nil
}

func (c *ConnInfo) Close() {
	_ = c.Client.Close()
}

func makePrivateKeySigner(privateKey []byte, passPhrase []byte) (gossh.Signer, error) {
	if len(passPhrase) != 0 {
		return gossh.ParsePrivateKeyWithPassphrase(privateKey, passPhrase)
	}
	return gossh.ParsePrivateKey(privateKey)
}
