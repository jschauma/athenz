/*
 *  Copyright The Athenz Authors
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package com.yahoo.athenz.common.server.msd.validator;

import com.yahoo.athenz.auth.PrivateKeyStore;
import com.yahoo.athenz.common.server.dns.HostnameResolver;
import com.yahoo.athenz.common.server.msd.MsdStore;
import com.yahoo.athenz.common.server.msd.repository.StaticWorkloadDataRepository;
import com.yahoo.athenz.msd.StaticWorkloadType;
import org.testng.annotations.Test;

import static org.testng.Assert.assertTrue;
import static org.testng.Assert.fail;

public class NoOpStaticWorkloadValidatorTest {

    @Test
    public void testNoOpValidator() {

        StaticWorkloadDataRepository<String> repository = new StaticWorkloadDataRepository<String>() {
            @Override
            public void initialize(PrivateKeyStore privateKeyStore, HostnameResolver hostnameResolver, MsdStore msdStore) {

            }

            @Override
            public String getDataByKey(String key) {
                return null;
            }
        };

        StaticWorkloadValidatorFactory factory = new StaticWorkloadValidatorFactory() {
            @Override
            public StaticWorkloadValidator create(StaticWorkloadType type, StaticWorkloadDataRepository<?> repository) {
                return StaticWorkloadValidatorFactory.super.create(type, repository);
            }
        };

        StaticWorkloadValidator validator = factory.create(StaticWorkloadType.VIP, repository);

        try {
            validator.initialize(repository);
            assertTrue(validator instanceof NoOpStaticWorkloadValidator);
            assertTrue(validator.validateWorkload("mydom", "mysvc", "xyz.athenz.io"));
        } catch (Exception ignored) {
            fail();
        }

    }
}