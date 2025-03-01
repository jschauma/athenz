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
import React from 'react';
import { render } from '@testing-library/react';
import GroupReviewPage from '../../../../../../pages/domain/[domain]/group/[group]/review';

describe('GroupReviewPage', () => {
    it('should render', () => {
        let domains = [];
        domains.push({ name: 'athens' });
        domains.push({ name: 'athens.ci' });
        let query = {
            domain: 'dom',
        };
        let domainDetails = {
            modified: '2020-02-12T21:44:37.792Z',
        };
        let groupDetails = {
            modified: '2020-02-12T21:44:37.792Z',
        };
        let headerDetails = {
            headerLinks: [
                {
                    title: 'Website',
                    url: 'http://www.athenz.io',
                    target: '_blank',
                },
            ],
            userData: {
                userLink: {
                    title: 'User Link',
                    url: '',
                    target: '_blank',
                },
            },
        };
        const { getByTestId } = render(
            <GroupReviewPage
                domains={domains}
                req='req'
                userId='userid'
                query={query}
                reload={false}
                domainDetails={domainDetails}
                groupDetails={groupDetails}
                domain='dom'
                domainResult={[]}
                headerDetails={headerDetails}
            />
        );
        const groupReviewPage = getByTestId('group-review');
        expect(groupReviewPage).toMatchSnapshot();
    });
});
