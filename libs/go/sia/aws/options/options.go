//
// Copyright Athenz Authors
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

package options

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/AthenZ/athenz/libs/go/sia/aws/doc"
	"github.com/AthenZ/athenz/libs/go/sia/aws/meta"
	"github.com/AthenZ/athenz/libs/go/sia/aws/stssession"
	"github.com/AthenZ/athenz/libs/go/sia/util"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

// package options contains types for parsing sia_config file and options to carry those config values

// ConfigService represents a service to be specified by user, and specify User/Group attributes for the service
type ConfigService struct {
	Filename       string `json:"filename,omitempty"`
	User           string `json:"user,omitempty"`
	Group          string `json:"group,omitempty"`
	ExpiryTime     int    `json:"expiry_time,omitempty"`
	SDSUdsUid      int    `json:"sds_uds_uid,omitempty"`
	SDSNodeId      string `json:"sds_node_id,omitempty"`
	SDSNodeCluster string `json:"sds_node_cluster,omitempty"`
}

// ConfigRole represents a role to be specified by user, and specify attributes for the role
type ConfigRole struct {
	Filename   string `json:"filename,omitempty"`
	ExpiryTime int    `json:"expiry_time,omitempty"`
}

// ConfigAccount represents each of the accounts that can be specified in the config file
type ConfigAccount struct {
	Name     string                `json:"name,omitempty"`     //name of the service identity
	User     string                `json:"user,omitempty"`     //the username to chown the cert/key dirs to. If absent, then root.
	Group    string                `json:"group,omitempty"`    //the group name to chown the cert/key dirs to. If absent, then athenz.
	Domain   string                `json:"domain,omitempty"`   //name of the domain for the identity
	Account  string                `json:"account,omitempty"`  //name of the account
	Service  string                `json:"service,omitempty"`  //name of the service for the identity
	Zts      string                `json:"zts,omitempty"`      //the ZTS to contact
	Filename string                `json:"filename,omitempty"` //filename to put the service certificate
	Roles    map[string]ConfigRole `json:"roles,omitempty"`    //map of roles to retrieve certificates for
	Version  string                `json:"version,omitempty"`  //sia version number
}

// Config represents entire sia_config file
type Config struct {
	Version         string                   `json:"version,omitempty"`           //name of the provider
	Service         string                   `json:"service,omitempty"`           //name of the service for the identity
	Services        map[string]ConfigService `json:"services,omitempty"`          //names of the multiple services for the identity
	Ssh             *bool                    `json:"ssh,omitempty"`               //ssh certificate support
	SanDnsWildcard  bool                     `json:"sandns_wildcard,omitempty"`   //san dns wildcard support
	UseRegionalSTS  bool                     `json:"regionalsts,omitempty"`       //whether to use a regional STS endpoint (default is false)
	Accounts        []ConfigAccount          `json:"accounts,omitempty"`          //array of configured accounts
	GenerateRoleKey bool                     `json:"generate_role_key,omitempty"` //private key to be generated for role certificate
	RotateKey       bool                     `json:"rotate_key,omitempty"`        //rotate private key support
	User            string                   `json:"user,omitempty"`              //the user name to chown the cert/key dirs to. If absent, then root
	Group           string                   `json:"group,omitempty"`             //the group name to chown the cert/key dirs to. If absent, then athenz
	SDSUdsPath      string                   `json:"sds_uds_path,omitempty"`      //uds path if the agent should support uds connections
	SDSUdsUid       int                      `json:"sds_uds_uid,omitempty"`       //uds connections must be from the given user uid
	ExpiryTime      int                      `json:"expiry_time,omitempty"`       //service and role certificate expiry in minutes
	RefreshInterval int                      `json:"refresh_interval,omitempty"`  //specifies refresh interval in minutes
	ZTSRegion       string                   `json:"zts_region,omitempty"`        //specifies zts region for the requests
}

// Role contains role details. Attributes are set based on the config values
type Role struct {
	Name     string
	Service  string
	Filename string
	User     string
	Uid      int
	Gid      int
	FileMode int
}

// Service represents service details. Attributes are filled in based on the config values
type Service struct {
	Name           string
	Filename       string
	User           string
	Group          string
	Uid            int
	Gid            int
	FileMode       int
	ExpiryTime     int
	SDSUdsUid      int
	SDSNodeId      string
	SDSNodeCluster string
}

// Options represents settings that are derived from config file and application defaults
type Options struct {
	Provider           string                //name of the provider
	Name               string                //name of the service identity
	User               string                //the user name to chown the cert/key dirs to. If absent, then root
	Group              string                //the group name to chown the cert/key dirs to. If absent, then athenz
	Domain             string                //name of the domain for the identity
	Account            string                //name of the account
	Service            string                //name of the service for the identity
	Zts                string                //the ZTS to contact
	Filename           string                //filename to put the service certificate
	InstanceId         string                //instance id if ec2, task id if running within eks/ecs
	Roles              map[string]ConfigRole //map of roles to retrieve certificates for
	Region             string                //region name
	SanDnsWildcard     bool                  //san dns wildcard support
	Version            string                //sia version number
	ZTSDomains         []string              //zts domain prefixes
	Services           []Service             //array of configured services
	Ssh                bool                  //ssh certificate support
	UseRegionalSTS     bool                  //use regional sts endpoint
	KeyDir             string                //private key directory path
	CertDir            string                //x.509 certificate directory path
	AthenzCACertFile   string                //filename to store Athenz CA certs
	ZTSCACertFile      string                //filename for CA certs when communicating with ZTS
	ZTSServerName      string                //ZTS server name, if necessary for tls
	ZTSAWSDomains      []string              //list of domain prefixes for sanDNS entries
	GenerateRoleKey    bool                  //option to generate a separate key for role certificates
	RotateKey          bool                  //rotate the private key when refreshing certificates
	BackUpDir          string                //backup directory for key/cert rotation
	CertCountryName    string                //generated x.509 certificate country name
	CertOrgName        string                //generated x.509 certificate organization name
	SshPubKeyFile      string                //ssh host public key file path
	SshCertFile        string                //ssh host certificate file path
	SshConfigFile      string                //sshd config file path
	PrivateIp          string                //instance private ip
	EC2Document        string                //EC2 instance identity document
	EC2Signature       string                //EC2 instance identity document pkcs7 signature
	EC2StartTime       *time.Time            //EC2 instance start time
	InstanceIdSanDNS   bool                  //include instance id in a san dns entry (backward compatible option)
	RolePrincipalEmail bool                  //include role principal in a san email field (backward compatible option)
	SDSUdsPath         string                //UDS path if the agent should support uds connections
	SDSUdsUid          int                   //UDS connections must be from the given user uid
	RefreshInterval    int                   //refresh interval for certificates - default 24 hours
	ZTSRegion          string                //ZTS region in case the client needs this information
}

func GetAccountId(metaEndPoint string, useRegionalSTS bool, region string) (string, error) {
	// first try to get the account from our creds and if
	// fails we'll fall back to identity document
	configAccount, err := InitCredsConfig("", useRegionalSTS, region)
	if err == nil {
		return configAccount.Account, nil
	}
	document, err := meta.GetData(metaEndPoint, "/latest/dynamic/instance-identity/document")
	if err != nil {
		return "", err
	}
	return doc.GetDocumentEntry(document, "accountId")
}

func InitCredsConfig(roleSuffix string, useRegionalSTS bool, region string) (*ConfigAccount, error) {
	account, domain, service, err := stssession.GetMetaDetailsFromCreds(roleSuffix, useRegionalSTS, region)
	if err != nil {
		return nil, fmt.Errorf("unable to parse role arn: %v", err)
	}
	return &ConfigAccount{
		Domain:  domain,
		Service: service,
		Account: account,
		Name:    fmt.Sprintf("%s.%s", domain, service),
	}, nil
}

func InitProfileConfig(metaEndPoint, roleSuffix string) (*ConfigAccount, error) {
	info, err := meta.GetData(metaEndPoint, "/latest/meta-data/iam/info")
	if err != nil {
		return nil, err
	}
	arn, err := doc.GetDocumentEntry(info, "InstanceProfileArn")
	if err != nil {
		return nil, err
	}
	account, domain, service, err := util.ParseRoleArn(arn, "instance-profile/", roleSuffix)
	if err != nil {
		return nil, err
	}

	return &ConfigAccount{
		Domain:  domain,
		Service: service,
		Account: account,
		Name:    fmt.Sprintf("%s.%s", domain, service),
	}, nil
}

func InitFileConfig(fileName, metaEndPoint string, useRegionalSTS bool, region, account string) (*Config, *ConfigAccount, error) {
	confBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, nil, err
	}
	if len(confBytes) == 0 {
		return nil, nil, errors.New("empty config bytes")
	}
	var config Config
	err = json.Unmarshal(confBytes, &config)
	if err != nil {
		return nil, nil, err
	}
	if config.Service == "" {
		return &config, nil, fmt.Errorf("missing required Service field from the config file")
	}
	// if we have more than one account block defined (not recommended)
	// then we need to determine our account id
	if len(config.Accounts) > 1 && account == "" {
		account, _ = GetAccountId(metaEndPoint, useRegionalSTS || config.UseRegionalSTS, region)
	}
	for _, configAccount := range config.Accounts {
		if configAccount.Account == account || len(config.Accounts) == 1 {
			if configAccount.Domain == "" || configAccount.Account == "" {
				return &config, nil, fmt.Errorf("missing required Domain/Account from the config file")
			}
			configAccount.Service = config.Service
			configAccount.Name = fmt.Sprintf("%s.%s", configAccount.Domain, configAccount.Service)
			return &config, &configAccount, nil
		}
	}
	return nil, nil, fmt.Errorf("missing account %s details from config file", account)
}

func InitEnvConfig(config *Config) (*Config, *ConfigAccount, error) {
	// it is possible that the config object was already created the
	// config file in which case we're not going to override any
	// of the settings.
	if config == nil {
		config = &Config{}
	}
	if !config.SanDnsWildcard {
		config.SanDnsWildcard = util.ParseEnvBooleanFlag("ATHENZ_SIA_SANDNS_WILDCARD")
	}
	if !config.UseRegionalSTS {
		config.UseRegionalSTS = util.ParseEnvBooleanFlag("ATHENZ_SIA_REGIONAL_STS")
	}
	if !config.GenerateRoleKey {
		config.GenerateRoleKey = util.ParseEnvBooleanFlag("ATHENZ_SIA_GENERATE_ROLE_KEY")
	}
	if !config.RotateKey {
		config.RotateKey = util.ParseEnvBooleanFlag("ATHENZ_SIA_ROTATE_KEY")
	}
	if config.User == "" {
		config.User = os.Getenv("ATHENZ_SIA_USER")
	}
	if config.Group == "" {
		config.Group = os.Getenv("ATHENZ_SIA_GROUP")
	}
	if config.SDSUdsPath == "" {
		config.SDSUdsPath = os.Getenv("ATHENZ_SIA_SDS_UDS_PATH")
	}
	if config.SDSUdsUid == 0 {
		uid := util.ParseEnvIntFlag("ATHENZ_SIA_SDS_UDS_UID", 0)
		if uid > 0 {
			config.SDSUdsUid = uid
		}
	}
	if config.ExpiryTime == 0 {
		expiryTime := util.ParseEnvIntFlag("ATHENZ_SIA_EXPIRY_TIME", 0)
		if expiryTime > 0 {
			config.ExpiryTime = expiryTime
		}
	}
	if config.RefreshInterval == 0 {
		refreshInterval := util.ParseEnvIntFlag("ATHENZ_SIA_REFRESH_INTERVAL", 0)
		if refreshInterval > 0 {
			config.RefreshInterval = refreshInterval
		}
	}
	if config.ZTSRegion == "" {
		config.ZTSRegion = os.Getenv("ATHENZ_SIA_ZTS_REGION")
	}

	roleArn := os.Getenv("ATHENZ_SIA_IAM_ROLE_ARN")
	if roleArn == "" {
		return config, nil, fmt.Errorf("athenz role arn env variable not configured")
	}
	account, domain, service, err := util.ParseRoleArn(roleArn, "role/", "")
	if err != nil {
		return config, nil, fmt.Errorf("unable to parse athenz role arn: %v", err)
	}
	if account == "" || domain == "" || service == "" {
		return config, nil, fmt.Errorf("invalid role arn - missing components: %s", roleArn)
	}
	return config, &ConfigAccount{
		Account: account,
		Domain:  domain,
		Service: service,
		Name:    fmt.Sprintf("%s.%s", domain, service),
	}, nil
}

// setOptions takes in sia_config objects and returns a pointer to Options after parsing and initializing the defaults
// It uses profile arn for defaults when sia_config is empty or non-parsable. It populates "services" array
func setOptions(config *Config, account *ConfigAccount, siaDir, version string) (*Options, error) {

	//update regional sts and wildcard settings based on config settings
	useRegionalSTS := false
	sanDnsWildcard := false
	sdsUdsPath := ""
	sdsUdsUid := 0
	generateRoleKey := false
	rotateKey := false
	expiryTime := 0
	refreshInterval := 24 * 60
	ztsRegion := ""

	if config != nil {
		useRegionalSTS = config.UseRegionalSTS
		sanDnsWildcard = config.SanDnsWildcard
		sdsUdsPath = config.SDSUdsPath
		sdsUdsUid = config.SDSUdsUid
		expiryTime = config.ExpiryTime
		ztsRegion = config.ZTSRegion
		if config.RefreshInterval > 0 {
			refreshInterval = config.RefreshInterval
		}

		//update account user/group settings if override provided at the config level
		if account.User == "" && config.User != "" {
			account.User = config.User
		}
		if account.Group == "" && config.Group != "" {
			account.Group = config.Group
		}

		//update generate role and rotate key options if config is provided
		generateRoleKey = config.GenerateRoleKey
		rotateKey = config.RotateKey
		if len(account.Roles) != 0 && generateRoleKey == false && rotateKey == true {
			log.Println("Cannot set rotate_key to true, with generate_role_key as false, when there are one or more roles defined in config")
			generateRoleKey = false
			rotateKey = false
		}
	}

	var services []Service
	if config == nil || len(config.Services) == 0 {
		// There is no sia_config, or multiple services are not configured. Populate services with the account information we gathered
		s := Service{
			Name:     account.Service,
			Filename: account.Filename,
			User:     account.User,
		}
		s.Uid, s.Gid, s.FileMode = util.SvcAttrs(account.User, account.Group)
		s.SDSUdsUid = sdsUdsUid
		s.ExpiryTime = expiryTime
		services = append(services, s)
	} else {
		// sia_config and services are found
		if _, ok := config.Services[config.Service]; !ok {
			return nil, fmt.Errorf("services: %+v mentioned, service: %q needs to be part of services", config.Services, config.Service)
		}
		// Populate config.Service into first
		first := Service{
			Name: config.Service,
		}

		// Populate the remaining into tail
		tail := []Service{}
		for name, s := range config.Services {
			svcExpiryTime := expiryTime
			if s.ExpiryTime > 0 {
				svcExpiryTime = s.ExpiryTime
			}
			svcSDSUdsUid := sdsUdsUid
			if s.SDSUdsUid != 0 {
				svcSDSUdsUid = s.SDSUdsUid
			}
			if name == config.Service {
				first.Filename = s.Filename
				first.User = s.User
				first.Group = s.Group
				// If User/Group are not specified, apply the User/Group settings from Config Account
				// This is for backwards compatibility - For other multiple services, the User/Group need to be explicitly mentioned in config
				if first.User == "" {
					first.User = account.User
				}
				if first.Group == "" {
					first.Group = account.Group
				}
				first.Uid, first.Gid, first.FileMode = util.SvcAttrs(first.User, first.Group)
				first.ExpiryTime = svcExpiryTime
				first.SDSNodeId = s.SDSNodeId
				first.SDSNodeCluster = s.SDSNodeCluster
				first.SDSUdsUid = svcSDSUdsUid
			} else {
				ts := Service{
					Name:     name,
					Filename: s.Filename,
					User:     s.User,
					Group:    s.Group,
				}
				ts.Uid, ts.Gid, ts.FileMode = util.SvcAttrs(s.User, s.Group)
				ts.ExpiryTime = svcExpiryTime
				ts.SDSNodeId = s.SDSNodeId
				ts.SDSNodeCluster = s.SDSNodeCluster
				ts.SDSUdsUid = svcSDSUdsUid
				tail = append(tail, ts)
			}
			if s.Filename != "" && s.Filename[0] == '/' {
				log.Println("when custom filepaths are specified, rotate_key and generate_role_key are not supported")
				generateRoleKey = false
				rotateKey = false
			}
		}
		services = append(services, first)
		services = append(services, tail...)
	}

	for _, r := range account.Roles {
		if r.Filename != "" && r.Filename[0] == '/' {
			log.Println("when custom filepaths are specified, rotate_key and generate_role_key are not supported")
			generateRoleKey = false
			rotateKey = false
			break
		}
	}

	return &Options{
		Name:             account.Name,
		User:             account.User,
		Group:            account.Group,
		Domain:           account.Domain,
		Account:          account.Account,
		Zts:              account.Zts,
		Filename:         account.Filename,
		Version:          fmt.Sprintf("SIA-AWS %s", version),
		UseRegionalSTS:   useRegionalSTS,
		SanDnsWildcard:   sanDnsWildcard,
		Services:         services,
		Roles:            account.Roles,
		CertDir:          fmt.Sprintf("%s/certs", siaDir),
		KeyDir:           fmt.Sprintf("%s/keys", siaDir),
		AthenzCACertFile: fmt.Sprintf("%s/certs/ca.cert.pem", siaDir),
		GenerateRoleKey:  generateRoleKey,
		RotateKey:        rotateKey,
		BackUpDir:        fmt.Sprintf("%s/backup", siaDir),
		SDSUdsPath:       sdsUdsPath,
		RefreshInterval:  refreshInterval,
		ZTSRegion:        ztsRegion,
	}, nil
}

// GetSvcNames returns comma separated list of service names
func GetSvcNames(svcs []Service) string {
	var b bytes.Buffer
	for _, svc := range svcs {
		b.WriteString(fmt.Sprintf("%s,", svc.Name))
	}
	return strings.TrimSuffix(b.String(), ",")
}

func NewOptions(config *Config, configAccount *ConfigAccount, siaDir, siaVersion string, useRegionalSTS bool, region string) (*Options, error) {

	opts, err := setOptions(config, configAccount, siaDir, siaVersion)
	if err != nil {
		log.Printf("Unable to formulate options, error: %v\n", err)
		return nil, err
	}

	opts.Region = region
	if useRegionalSTS {
		opts.UseRegionalSTS = true
	}
	return opts, nil
}
