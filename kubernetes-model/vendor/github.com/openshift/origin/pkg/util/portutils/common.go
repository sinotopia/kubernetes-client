/**
 * Copyright (C) 2015 Red Hat, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package portutils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fsouza/go-dockerclient"
)

// SplitPortAndProtocol splits a string of format port/proto and returns
// a docker.Port
func SplitPortAndProtocol(port string) (docker.Port, error) {
	dp := docker.Port(port)

	err := ValidatePortAndProtocol(dp)

	return dp, err
}

// SplitPortAndProtocolArray splits an array of strings of format port/proto
// and returns an array of docker.Port
func SplitPortAndProtocolArray(ports []string) ([]docker.Port, []error) {
	allErrs := []error{}
	allPorts := []docker.Port{}

	for _, port := range ports {
		dp, err := SplitPortAndProtocol(port)

		if err != nil {
			allErrs = append(allErrs, err)
		}

		allPorts = append(allPorts, dp)
	}

	if len(allErrs) > 0 {
		return allPorts, allErrs
	}

	return allPorts, nil
}

// ValidatePortAndProtocol validates the port range and protocol of a docker.Port
func ValidatePortAndProtocol(port docker.Port) error {
	errs := []string{}

	_, err := strconv.ParseUint(port.Port(), 10, 16)
	if err != nil {
		if numError, ok := err.(*strconv.NumError); ok {
			if numError.Err == strconv.ErrRange || numError.Err == strconv.ErrSyntax {
				errs = append(errs, "port number must be in range 0 - 65535")
			}
		}
	}

	if len(port.Proto()) > 0 && !(strings.ToUpper(port.Proto()) == "TCP" || strings.ToUpper(port.Proto()) == "UDP") {
		errs = append(errs, "protocol must be tcp or udp")
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to parse port %s/%s: [%v]", port.Port(), port.Proto(), strings.Join(errs, ", "))
	}

	return nil
}