package config

import "github.com/sdslabs/pinger/pkg/database"

// DBConn config
type DBConn struct {
	Host     string
	Port     uint16
	Name     string
	Username string
	Password string
	SSLMode  bool
}

// GetHost returns the host of the database provider.
func (m *DBConn) GetHost() string {
	return m.Host
}

// GetPort returns the port of the database provider.
func (m *DBConn) GetPort() uint16 {
	return m.Port
}

// GetName returns the name of the provider.
func (m *DBConn) GetName() string {
	return m.Name
}

// GetUsername returns the username of the database provider.
func (m *DBConn) GetUsername() string {
	return m.Username
}

// GetPassword returns the password of the database provider.
func (m *DBConn) GetPassword() string {
	return m.Password
}

// IsSSLMode tells if the connection with the provider is to be established
// through SSL.
func (m *DBConn) IsSSLMode() bool {
	return m.SSLMode
}

// Interface guard
var _ database.Config = (*DBConn)(nil)
