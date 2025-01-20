//go:build !windows
// +build !windows

package workloadapi

import "errors"

// appendDialOptionsOS appends OS specific dial options
func (c *Client) appendDialOptionsOS() {
	// No options to add in this platform
}
func (c *Client) setAddress() error {
	if c.config.namedPipeName != "" {
		// Purely defensive. This should never happen.
		return errors.New("named pipes not supported in this platform")
	}

	if c.config.address == "" {
		var ok bool
		c.config.address, ok = GetDefaultAddress()
		if !ok {
			return errors.New("workload endpoint socket address is not configured")
		}
	}

	var err error
<<<<<<< HEAD
	c.config.address, err = TargetFromAddress(c.config.address)
=======
	c.config.address, err = parseTargetFromStringAddr(c.config.address)
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	return err
}
