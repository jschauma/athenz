/*
 * Copyright The Athenz Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package com.yahoo.athenz.instance.provider.impl;

import com.yahoo.athenz.auth.KeyStore;
import com.yahoo.athenz.auth.util.Crypto;
import com.yahoo.athenz.common.server.dns.HostnameResolver;
import com.yahoo.athenz.instance.provider.*;
import org.bouncycastle.asn1.x500.style.BCStyle;
import org.eclipse.jetty.util.StringUtil;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import javax.net.ssl.SSLContext;
import java.util.*;

public class SecureBootProvider implements InstanceProvider {

    private static final Logger LOG = LoggerFactory.getLogger(SecureBootProvider.class);
    private static final String URI_HOSTNAME_PREFIX = "athenz://hostname/";

    static final String ZTS_PROP_SB_PROVIDER_DNS_SUFFIX  = "athenz.zts.sb_provider_dns_suffix";
    static final String ZTS_PROP_SB_PRINCIPAL_LIST       = "athenz.zts.sb_provider_service_list";
    static final String ZTS_PROP_SB_ISSUER_DN_LIST       = "athenz.zts.sb_provider_issuer_dn_list";
    static final String ZTS_PROP_SB_ATTR_VALIDATOR_FACTORY_CLASS = "athenz.zts.sb_provider_attr_validator_factory_class";
    static final String ZTS_SB_PROVIDER_SERVICE  = "sys.auth.sb-provider";

    // To look up any secrets needed by the provider
    KeyStore keyStore = null;
    Set<String> dnsSuffixes = null;
    String provider = null;
    Set<String> principals = null;
    Set<String> issuerDNs = null;
    HostnameResolver hostnameResolver = null;
    AttrValidator attrValidator = null;

    @Override
    public Scheme getProviderScheme() {
        return Scheme.CLASS;
    }

    @Override
    public void initialize(String provider, String providerEndpoint, SSLContext sslContext,
            KeyStore keyStore) {

        // save our provider name

        this.provider = provider;

        // obtain list of valid principals for this principal if one is configured

        final String principalList = System.getProperty(ZTS_PROP_SB_PRINCIPAL_LIST);
        if (principalList != null && !principalList.isEmpty()) {
            principals = new HashSet<>(Arrays.asList(principalList.split(",")));
        }

        final String issuerList = System.getProperty(ZTS_PROP_SB_ISSUER_DN_LIST);
        if (issuerList != null && !issuerList.isEmpty()) {
            issuerDNs = new HashSet<>(Arrays.asList(issuerList.split(";")));
        }

        // determine the dns suffix. if this is not specified we'll just default to zts.athenz.cloud

        dnsSuffixes = new HashSet<>();
        String dnsSuffix = System.getProperty(ZTS_PROP_SB_PROVIDER_DNS_SUFFIX, "zts.athenz.cloud");
        dnsSuffixes.addAll(Arrays.asList(dnsSuffix.split(",")));

        this.keyStore = keyStore;
        this.attrValidator = newAttrValidator();
    }

    static AttrValidator newAttrValidator() {
        final String factoryClass = System.getProperty(ZTS_PROP_SB_ATTR_VALIDATOR_FACTORY_CLASS);
        if (factoryClass == null) {
            return null;
        }

        AttrValidatorFactory attrValidatorFactory;
        try {
            attrValidatorFactory = (AttrValidatorFactory) Class.forName(factoryClass).getConstructor().newInstance();
        } catch (Exception e) {
            LOG.error("Invalid AttributeValidatorFactory class: {}", factoryClass, e);
            throw new IllegalArgumentException("Invalid AttributeValidatorFactory class");
        }

        return attrValidatorFactory.create();
    }

    @Override
    public void setHostnameResolver(HostnameResolver hostnameResolver) {
        this.hostnameResolver = hostnameResolver;
    }

    public void setAttrValidator(AttrValidator attrValidator) {
        this.attrValidator = attrValidator;
    }

    private ResourceException forbiddenError(String message, String txt) {
        LOG.error("mesaage: {}, logText: {}", message, txt);
        return new ResourceException(ResourceException.FORBIDDEN, message);
    }

    @Override
    public InstanceConfirmation confirmInstance(InstanceConfirmation confirmation) {
        return validateInstanceRequest(confirmation, true);
    }

    @Override
    public InstanceConfirmation refreshInstance(InstanceConfirmation confirmation) {
        return validateInstanceRequest(confirmation, false);
    }

    InstanceConfirmation validateInstanceRequest(InstanceConfirmation confirmation, final boolean register) {

        final String instanceDomain = confirmation.getDomain();
        final String instanceService = confirmation.getService();

        final Map<String, String> instanceAttributes = confirmation.getAttributes();

        // make sure this service has been configured to be supported
        // by this provider

        if (principals != null && !principals.contains(instanceDomain + "." + instanceService)) {
            throw forbiddenError("Service not supported to be launched by SecureBoot Provider",
                    logTxt(confirmation));
        }

        final String attestationData = confirmation.getAttestationData();
        if (StringUtil.isEmpty(attestationData)) {
            throw forbiddenError("Invalid attestation data", logTxt(confirmation));
        }

        // make sure the provider matches
        if (!ZTS_SB_PROVIDER_SERVICE.equals(confirmation.getProvider())) {
            throw forbiddenError("Invalid provider in the confirmation", logTxt(confirmation));
        }

        // for register requests, validate that the request is from a known issuer
        if (register && !validateIssuer(instanceAttributes)) {
            throw forbiddenError("Invalid issuer", logTxt(confirmation));
        }

        final String hostname = InstanceUtils.getInstanceProperty(instanceAttributes,
                InstanceProvider.ZTS_INSTANCE_HOSTNAME);

        // validate the hostname (from request cert CN for register or from URI for refresh)
        if (!validateHostname(hostname, register, instanceAttributes)) {
            throw forbiddenError("Unable to validate certificate request hostname", logTxt(confirmation));
        }

        // Confirm the instance attributes as per the attribute validator
        if (!attrValidator.confirm(confirmation)) {
            throw forbiddenError("Unable to validate request instance attributes", logTxt(confirmation));
        }

        // validate SanIP, if available
        if (!validateSanIp(hostname, instanceAttributes)) {
            throw forbiddenError("Unable to validate request IP address", logTxt(confirmation));
        }

        // validate the certificate san DNS names
        StringBuilder instanceId = new StringBuilder(256);
        if (!InstanceUtils.validateCertRequestSanDnsNames(instanceAttributes, instanceDomain,
                instanceService, dnsSuffixes, instanceId)) {
            throw forbiddenError("Unable to validate certificate request DNS", logTxt(confirmation));
        }

        Map<String, String> attributes = new HashMap<>();
        confirmation.setAttributes(attributes);
        return confirmation;
    }

    /**
     * validateIssuer ensures that IssuerDN is passed in.
     * If issuerDNs is configured, it validates the IssuerDN against the configured values
     * @param attributes
     * @return
     */
    boolean validateIssuer(final Map<String, String> attributes) {
        final String dn = InstanceUtils.getInstanceProperty(attributes,
                InstanceProvider.ZTS_INSTANCE_CERT_ISSUER_DN);

        if (StringUtil.isEmpty(dn)) {
            LOG.error("issuer DN must be passed by ZTS");
            return false;
        }

        // If no allow list configured, accept any issuer
        if (issuerDNs == null) {
            return true;
        }

        return issuerDNs.contains(dn);
    }

    /**
     * validateHostname validates the hostname passed in for the instance
     * @param hostname of the instance
     * @param register whether the request is for register or refresh
     * @param attributes from the confirmation
     * @return true or false
     */
    static boolean validateHostname(final String hostname, final boolean register, final Map<String, String> attributes) {
        if (StringUtil.isEmpty(hostname)) {
            return false;
        }

        return register ? validateCnHostname(hostname, attributes)
                : validateCertHostname(hostname, attributes);
    }

    static String getSubjectCn(String dn) {
        return Crypto.extractX500DnField(dn, BCStyle.CN);
    }

    static boolean validateCnHostname(final String hostname, final Map<String, String> attributes) {
        final String subjectDn = InstanceUtils.getInstanceProperty(attributes,
                InstanceProvider.ZTS_INSTANCE_CERT_SUBJECT_DN);
        final String cn = getSubjectCn(subjectDn);
        return hostname.equals(cn);
    }

    /**
     * validateCertHostname matches the instance hostname against the value passed in
     * by ZTS in cert hostname, which ZTS extracts it from SAN URI
     * @param hostname
     * @param attributes
     * @return
     */
    static boolean validateCertHostname(final String hostname, final Map<String, String> attributes) {
        final String certHostname = InstanceUtils.getInstanceProperty(attributes,
                InstanceProvider.ZTS_INSTANCE_CERT_HOSTNAME);

        // It's possible the refresh request is coming in without the hostname in the URI
        if (certHostname == null || certHostname.isEmpty()) {
            return true;
        }
        return hostname.equals(certHostname);
    }

    /**
     * verifies that the IPs in SanIP map to the hostname via the hostIps fetched from the hostnameResolver
     * @param hostname name of the instance
     * @param attributes passsed in from the instance confirmation
     * @return true if sanIps is null/empty or all sanIps match hostIps. false otherwise
     */
    boolean validateSanIp(final String hostname, final Map<String, String> attributes) {

        final String sanIpStr = InstanceUtils.getInstanceProperty(attributes,
                InstanceProvider.ZTS_INSTANCE_SAN_IP);

        String[] sanIps = null;
        if (sanIpStr != null && !sanIpStr.isEmpty()) {
            sanIps = sanIpStr.split(",");
        }

        if (sanIps == null || sanIps.length == 0) {
            return true;
        }

        Set<String> hostIps = hostnameResolver.getAllByName(hostname);;

        LOG.debug("Validating sanIps: {}, hostIps: {}", sanIps, hostIps);

        for (String sanIp: sanIps) {
            if (!hostIps.contains(sanIp)) {
                LOG.error("Unable to match sanIp: {} with hostIps:{}", sanIps, hostIps);
                return false;
            }
        }

        return true;
    }

    public static String logTxt(InstanceConfirmation confirmation) {
        return "InstanceConfirmation{" +
                "provider='" + confirmation.getProvider() + '\'' +
                ", domain='" + confirmation.getDomain() + '\'' +
                ", service='" + confirmation.getService() + '\'' +
                '}';
    }
}
