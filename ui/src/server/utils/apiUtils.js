/*
 * Copyright The Athenz Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */
'use strict';
const userService = require('../services/userService');

module.exports.getPendingDomainMemberData = (values) => {
    let data = values[1];
    let pendingMap = {};
    data.domainRoleMembersList.forEach((domain) => {
        const domainName = domain.domainName;
        domain.members.forEach((member) => {
            const memberName = member.memberName;
            member.memberRoles.forEach((role) => {
                const roleName = role.roleName;
                const expiryDate = role.expiration || null;
                const userComment = role.auditRef || null;
                const key = domainName + memberName + roleName;
                pendingMap[key] = {
                    category: 'role',
                    domainName: domainName,
                    memberName: memberName,
                    memberNameFull: userService.getUserFullName(memberName),
                    roleName: roleName,
                    userComment: userComment,
                    auditRef: '',
                    requestPrincipal: role.requestPrincipal,
                    requestPrincipalFull: userService.getUserFullName(
                        role.requestPrincipal
                    ),
                    requestTime: role.requestTime,
                    expiryDate: expiryDate,
                };
            });
        });
    });
    values[0].domainGroupMembersList.forEach((domain) => {
        const domainName = domain.domainName;
        domain.members.forEach((member) => {
            const memberName = member.memberName;
            member.memberGroups.forEach((group) => {
                const groupName = group.groupName;
                const expiryDate = group.expiration || null;
                const userComment = group.auditRef || null;
                const key = domainName + memberName + groupName;
                pendingMap[key] = {
                    category: 'group',
                    domainName: domainName,
                    memberName: memberName,
                    memberNameFull: userService.getUserFullName(memberName),
                    roleName: groupName,
                    userComment: userComment,
                    auditRef: '',
                    requestPrincipal: group.requestPrincipal,
                    requestPrincipalFull: userService.getUserFullName(
                        group.requestPrincipal
                    ),
                    requestTime: group.requestTime,
                    expiryDate: expiryDate,
                };
            });
        });
    });
    return pendingMap;
};

module.exports.getPendingDomainMembersPromise = (params, req) => {
    let promises = [];
    promises.push(
        new Promise((resolve, reject) => {
            req.clients.zms.getPendingDomainGroupMembersList(
                params,
                (err, data) => {
                    if (err) {
                        reject(err);
                    }
                    if (data) {
                        resolve(data);
                    }
                }
            );
        })
    );
    promises.push(
        new Promise((resolve, reject) => {
            req.clients.zms.getPendingDomainRoleMembersList(
                params,
                (err, data) => {
                    if (err) {
                        reject(err);
                    }
                    if (data) {
                        resolve(data);
                    }
                }
            );
        })
    );
    return promises;
};
