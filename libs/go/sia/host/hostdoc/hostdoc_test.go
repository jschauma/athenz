//
// Copyright The Athenz Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package hostdoc

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const HOSTDOC_STR = `
{
   "domain": "sports",
   "service": "soccer",
   "profile": "prod",
   "uuid": "3e4c2da84a264d718b218ce58b1b3b8f",
   "zone": "west"
}
`

const HOSTDOC_STR_SVCS = `
{
   "domain": "sports",
   "services": "soccer,nfl,yamas,splunk,logger",
   "uuid": "3e4c2da84a264d718b218ce58b1b3b8f",
   "zone": "west"
}
`

// Todo: Improve the tests here to be able to parse both "service" and "services"

func TestNewHostDocParseErr(t *testing.T) {
	a := assert.New(t)

	_, _, err := NewPlainDoc([]byte(`{`))
	a.NotNil(err)
	a.Containsf(err.Error(), "unexpected end of JSON input", "unexpected err: %v", err)
}

func TestNewHostDoc(t *testing.T) {
	a := assert.New(t)

	hostDoc, _, err := NewPlainDoc([]byte(HOSTDOC_STR))
	a.Nil(err)
	a.Equal("sports", hostDoc.Domain, "domain should match")
	a.Equal("soccer", hostDoc.Services[0], "service should match")
	a.Equal("prod", hostDoc.Profile, "profile should match")
	a.Equal("3e4c2da8-4a26-4d71-8b21-8ce58b1b3b8f", hostDoc.Uuid, "Uuid should match")
	a.Equal("west", hostDoc.Zone, "service should match")
}

func TestNewHostDocServices(t *testing.T) {
	a := assert.New(t)

	hostDoc, _, err := NewPlainDoc([]byte(HOSTDOC_STR_SVCS))
	a.Nil(err)
	a.Equal("sports", hostDoc.Domain, "domain should match")
	a.Equal(hostDoc.Domain, "sports", "domain should match")
	a.Equal("soccer", hostDoc.Services[0], "service should match")
	a.Equal("nfl", hostDoc.Services[1], "service should match")
	a.Equal("logger", hostDoc.Services[4], "service should match")
	a.Equal("", hostDoc.Profile, "profile should be empty")
	a.Equal("3e4c2da8-4a26-4d71-8b21-8ce58b1b3b8f", hostDoc.Uuid, "Uuid should match")
	a.Equal("west", hostDoc.Zone, "service should match")
}

func TestNewHostDocNoDomain(t *testing.T) {
	a := assert.New(t)

	docBytes, err := ioutil.ReadFile(filepath.Join("testdata", "host_document.nodomain"))
	require.Nilf(t, err, "unexpected err: %v", err)

	_, _,  err = NewPlainDoc(docBytes)
	a.NotNil(err)
	a.Containsf(err.Error(), "unable to find domain in host_document,", "unexpected err: %v", err)
}

func TestNewHostDocMultipleSvcs(t *testing.T) {
	a := assert.New(t)

	docBytes, err := ioutil.ReadFile(filepath.Join("testdata", "host_document.services"))
	require.Nilf(t, err, "unexpected err: %v", err)

	hostDoc, pvdrStr,  err := NewPlainDoc(docBytes)
	a.Nil(err)
	a.Equal(hostDoc.Domain, "athenz.examples", "domain should match")
	a.Equalf("httpd", hostDoc.Services[0], "unexpected service, svcs: %+v", hostDoc.Services)
	a.Equalf("ftpd", hostDoc.Services[1], "unexpected service, svcs: %+v", hostDoc.Services)
	a.Equalf("infra.provider-baremetal", pvdrStr, "unexpected provider, provider: %+v", hostDoc.Provider)
}

func TestNewHostDocNoService(t *testing.T) {
	a := assert.New(t)

	docBytes, err := ioutil.ReadFile(filepath.Join("testdata", "host_document.noservice"))
	require.Nilf(t, err, "unexpected err: %v", err)

	_, _,  err = NewPlainDoc(docBytes)
	a.NotNil(err)
	a.Containsf(err.Error(), "unable to find service or services in host_document,", "unexpected err: %v", err)
}

func TestNewHostDocIp(t *testing.T) {
	a := assert.New(t)

	docBytes, err := ioutil.ReadFile(filepath.Join("testdata", "host_document.ip"))
	require.Nilf(t, err, "unexpected err: %v", err)

	hostDoc, pvdrStr,  err := NewPlainDoc(docBytes)
	a.Nil(err)
	a.Equal(hostDoc.Domain, "athenz.examples", "domain should match")
	a.Equalf("infra.provider-baremetal", pvdrStr, "unexpected provider, provider: %+v", hostDoc.Provider)
	a.Truef(hostDoc.Ip["10.196.66.191"], "unexpected ips: %+v", hostDoc.Ip)
	a.Truef(hostDoc.Ip["2a00:1288:110:c30::1027"], "unexpected ips: %+v", hostDoc.Ip)
}

func TestNewHostDocIpUncompressed(t *testing.T) {
	a := assert.New(t)

	docBytes, err := ioutil.ReadFile(filepath.Join("testdata", "host_document.ip.uncompressed"))
	require.Nilf(t, err, "unexpected err: %v", err)

	hostDoc, pvdrStr,  err := NewPlainDoc(docBytes)
	a.Nil(err)
	a.Equal(hostDoc.Domain, "athenz.examples", "domain should match")
	a.Equalf("infra.provider-baremetal", pvdrStr, "unexpected provider, provider: %+v", hostDoc.Provider)
	a.Lenf(hostDoc.Ip, 2, "unexpected number of ips: %+v", hostDoc.Ip)
	a.Truef(hostDoc.Ip["10.196.66.191"], "unexpected ips: %+v", hostDoc.Ip)
	a.Truef(hostDoc.Ip["2a00:1288:110:c30::1027"], "unexpected ips: %+v", hostDoc.Ip)
}

func TestNewHostDocIpUpperCase(t *testing.T) {
	a := assert.New(t)

	docBytes, err := ioutil.ReadFile(filepath.Join("testdata", "host_document.ip.uppercase"))
	require.Nilf(t, err, "unexpected err: %v", err)

	hostDoc, pvdrStr,  err := NewPlainDoc(docBytes)
	a.Nil(err)
	a.Equal(hostDoc.Domain, "athenz.examples", "domain should match")
	a.Equalf("infra.provider-baremetal", pvdrStr, "unexpected provider, provider: %+v", hostDoc.Provider)
	a.Truef(hostDoc.Ip["216.115.99.107"], "unexpected ips: %+v", hostDoc.Ip)
	a.Truef(hostDoc.Ip["2001:4998:f00a:1fc::5301"], "unexpected ips: %+v", hostDoc.Ip)
}

