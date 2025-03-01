/*
 * Copyright 2017 Yahoo Inc.
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

package com.yahoo.athenz.zts;

import org.testng.annotations.Test;

import java.util.Collections;

import static org.testng.Assert.*;

public class InstanceRefreshInformationTest {

    @Test
    public void testInstanceRefreshInformation() {

        InstanceRefreshInformation i1 = new InstanceRefreshInformation();
        InstanceRefreshInformation i2 = new InstanceRefreshInformation();

        // set
        i1.setAttestationData("doc");
        i1.setCsr("sample_csr");
        i1.setSsh("ssh");
        i1.setToken(false);
        i1.setExpiryTime(180);
        i1.setHostname("host1.athenz.cloud");
        i1.setHostCnames(Collections.singletonList("host1"));
        i1.setSshCertRequest(new SSHCertRequest());

        i2.setAttestationData("doc");
        i2.setCsr("sample_csr");
        i2.setSsh("ssh");
        i2.setToken(false);
        i2.setExpiryTime(180);
        i2.setHostname("host1.athenz.cloud");
        i2.setHostCnames(Collections.singletonList("host1"));
        i2.setSshCertRequest(new SSHCertRequest());

        // getter assertion
        assertEquals(i1.getAttestationData(), "doc");
        assertEquals(i1.getCsr(), "sample_csr");
        assertEquals(i1.getSsh(), "ssh");
        assertEquals(i1.getToken(), Boolean.FALSE);
        assertEquals(i1.getExpiryTime(), Integer.valueOf(180));
        assertEquals(i1.getHostname(), "host1.athenz.cloud");
        assertEquals(i1.getHostCnames(), Collections.singletonList("host1"));
        assertEquals(i1.getSshCertRequest(), new SSHCertRequest());

        assertEquals(i2, i1);
        assertEquals(i2, i2);
        assertNotEquals(i2, null);
        assertNotEquals(i2, "string");
        assertNotEquals("string", i2);

        i2.setAttestationData(null);
        assertNotEquals(i1, i2);
        i2.setAttestationData("doc2");
        assertNotEquals(i1, i2);
        i2.setAttestationData("doc");
        assertEquals(i1, i2);

        i2.setCsr(null);
        assertNotEquals(i1, i2);
        i2.setCsr("csr2");
        assertNotEquals(i1, i2);
        i2.setCsr("sample_csr");
        assertEquals(i1, i2);

        i2.setSsh(null);
        assertNotEquals(i1, i2);
        i2.setSsh("ssh2");
        assertNotEquals(i1, i2);
        i2.setSsh("ssh");
        assertEquals(i1, i2);

        i2.setToken(null);
        assertNotEquals(i1, i2);
        i2.setToken(true);
        assertNotEquals(i1, i2);
        i2.setToken(false);
        assertEquals(i1, i2);

        i2.setExpiryTime(null);
        assertNotEquals(i1, i2);
        i2.setExpiryTime(120);
        assertNotEquals(i1, i2);
        i2.setExpiryTime(180);
        assertEquals(i1, i2);

        i2.setHostname(null);
        assertNotEquals(i1, i2);
        i2.setHostname("host2");
        assertNotEquals(i1, i2);
        i2.setHostname("host1.athenz.cloud");
        assertEquals(i1, i2);

        i2.setHostCnames(null);
        assertNotEquals(i1, i2);
        i2.setHostCnames(Collections.singletonList("host1-cname1.athenz.cloud"));
        assertNotEquals(i1, i2);
        i2.setHostCnames(Collections.singletonList("host1"));
        assertEquals(i1, i2);

        i2.setSshCertRequest(null);
        assertNotEquals(i1, i2);
        i2.setSshCertRequest(new SSHCertRequest().setCsr("csr"));
        assertNotEquals(i1, i2);
        i2.setSshCertRequest(new SSHCertRequest());
        assertEquals(i1, i2);
    }
}
