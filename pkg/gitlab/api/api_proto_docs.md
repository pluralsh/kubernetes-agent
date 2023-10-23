# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [pkg/gitlab/api/api.proto](#pkg_gitlab_api_api-proto)
    - [AccessAsAgentAuthorization](#gitlab-agent-gitlab-api-AccessAsAgentAuthorization)
    - [AccessAsProxyAuthorization](#gitlab-agent-gitlab-api-AccessAsProxyAuthorization)
    - [AccessAsUserAuthorization](#gitlab-agent-gitlab-api-AccessAsUserAuthorization)
    - [AgentConfigurationRequest](#gitlab-agent-gitlab-api-AgentConfigurationRequest)
    - [AllowedAgent](#gitlab-agent-gitlab-api-AllowedAgent)
    - [AllowedAgentsForJob](#gitlab-agent-gitlab-api-AllowedAgentsForJob)
    - [AuthorizeProxyUserRequest](#gitlab-agent-gitlab-api-AuthorizeProxyUserRequest)
    - [AuthorizeProxyUserResponse](#gitlab-agent-gitlab-api-AuthorizeProxyUserResponse)
    - [AuthorizedAgentForUser](#gitlab-agent-gitlab-api-AuthorizedAgentForUser)
    - [ConfigProject](#gitlab-agent-gitlab-api-ConfigProject)
    - [Configuration](#gitlab-agent-gitlab-api-Configuration)
    - [Environment](#gitlab-agent-gitlab-api-Environment)
    - [GetAgentInfoResponse](#gitlab-agent-gitlab-api-GetAgentInfoResponse)
    - [GetProjectInfoResponse](#gitlab-agent-gitlab-api-GetProjectInfoResponse)
    - [Group](#gitlab-agent-gitlab-api-Group)
    - [GroupAccessCF](#gitlab-agent-gitlab-api-GroupAccessCF)
    - [Job](#gitlab-agent-gitlab-api-Job)
    - [Pipeline](#gitlab-agent-gitlab-api-Pipeline)
    - [Project](#gitlab-agent-gitlab-api-Project)
    - [ProjectAccessCF](#gitlab-agent-gitlab-api-ProjectAccessCF)
    - [User](#gitlab-agent-gitlab-api-User)
  
- [Scalar Value Types](#scalar-value-types)



<a name="pkg_gitlab_api_api-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## pkg/gitlab/api/api.proto



<a name="gitlab-agent-gitlab-api-AccessAsAgentAuthorization"></a>

### AccessAsAgentAuthorization







<a name="gitlab-agent-gitlab-api-AccessAsProxyAuthorization"></a>

### AccessAsProxyAuthorization



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| agent | [AccessAsAgentAuthorization](#gitlab-agent-gitlab-api-AccessAsAgentAuthorization) |  |  |
| user | [AccessAsUserAuthorization](#gitlab-agent-gitlab-api-AccessAsUserAuthorization) |  |  |






<a name="gitlab-agent-gitlab-api-AccessAsUserAuthorization"></a>

### AccessAsUserAuthorization



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| projects | [ProjectAccessCF](#gitlab-agent-gitlab-api-ProjectAccessCF) | repeated |  |
| groups | [GroupAccessCF](#gitlab-agent-gitlab-api-GroupAccessCF) | repeated |  |






<a name="gitlab-agent-gitlab-api-AgentConfigurationRequest"></a>

### AgentConfigurationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| agent_id | [int64](#int64) |  |  |
| agent_config | [gitlab.agent.agentcfg.ConfigurationFile](#gitlab-agent-agentcfg-ConfigurationFile) |  |  |






<a name="gitlab-agent-gitlab-api-AllowedAgent"></a>

### AllowedAgent



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [int64](#int64) |  |  |
| config_project | [ConfigProject](#gitlab-agent-gitlab-api-ConfigProject) |  |  |
| configuration | [Configuration](#gitlab-agent-gitlab-api-Configuration) |  |  |






<a name="gitlab-agent-gitlab-api-AllowedAgentsForJob"></a>

### AllowedAgentsForJob



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| allowed_agents | [AllowedAgent](#gitlab-agent-gitlab-api-AllowedAgent) | repeated |  |
| job | [Job](#gitlab-agent-gitlab-api-Job) |  |  |
| pipeline | [Pipeline](#gitlab-agent-gitlab-api-Pipeline) |  |  |
| project | [Project](#gitlab-agent-gitlab-api-Project) |  |  |
| user | [User](#gitlab-agent-gitlab-api-User) |  |  |
| environment | [Environment](#gitlab-agent-gitlab-api-Environment) |  | optional |






<a name="gitlab-agent-gitlab-api-AuthorizeProxyUserRequest"></a>

### AuthorizeProxyUserRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| agent_id | [int64](#int64) |  |  |
| access_type | [string](#string) |  |  |
| access_key | [string](#string) |  |  |
| csrf_token | [string](#string) |  |  |






<a name="gitlab-agent-gitlab-api-AuthorizeProxyUserResponse"></a>

### AuthorizeProxyUserResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| agent | [AuthorizedAgentForUser](#gitlab-agent-gitlab-api-AuthorizedAgentForUser) |  |  |
| user | [User](#gitlab-agent-gitlab-api-User) |  |  |
| access_as | [AccessAsProxyAuthorization](#gitlab-agent-gitlab-api-AccessAsProxyAuthorization) |  |  |






<a name="gitlab-agent-gitlab-api-AuthorizedAgentForUser"></a>

### AuthorizedAgentForUser



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [int64](#int64) |  |  |
| config_project | [ConfigProject](#gitlab-agent-gitlab-api-ConfigProject) |  |  |






<a name="gitlab-agent-gitlab-api-ConfigProject"></a>

### ConfigProject



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [int64](#int64) |  |  |






<a name="gitlab-agent-gitlab-api-Configuration"></a>

### Configuration
Configuration contains shared fields from agentcfg.CiAccessProjectCF and agentcfg.CiAccessGroupCF.
It is used to parse response from the allowed_agents API endpoint.
See https://gitlab.com/gitlab-org/cluster-integration/gitlab-agent/-/blob/master/doc/kubernetes_ci_access.md#apiv4joballowed_agents-api.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| default_namespace | [string](#string) |  |  |
| access_as | [gitlab.agent.agentcfg.CiAccessAsCF](#gitlab-agent-agentcfg-CiAccessAsCF) |  |  |






<a name="gitlab-agent-gitlab-api-Environment"></a>

### Environment



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| slug | [string](#string) |  |  |
| tier | [string](#string) |  |  |






<a name="gitlab-agent-gitlab-api-GetAgentInfoResponse"></a>

### GetAgentInfoResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [int64](#int64) |  |  |
| agent_id | [int64](#int64) |  |  |
| agent_name | [string](#string) |  |  |
| default_branch | [string](#string) |  |  |






<a name="gitlab-agent-gitlab-api-GetProjectInfoResponse"></a>

### GetProjectInfoResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [int64](#int64) |  |  |
| default_branch | [string](#string) |  |  |






<a name="gitlab-agent-gitlab-api-Group"></a>

### Group



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [int64](#int64) |  |  |






<a name="gitlab-agent-gitlab-api-GroupAccessCF"></a>

### GroupAccessCF



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [int64](#int64) |  |  |
| roles | [string](#string) | repeated |  |






<a name="gitlab-agent-gitlab-api-Job"></a>

### Job



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [int64](#int64) |  |  |






<a name="gitlab-agent-gitlab-api-Pipeline"></a>

### Pipeline



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [int64](#int64) |  |  |






<a name="gitlab-agent-gitlab-api-Project"></a>

### Project



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [int64](#int64) |  |  |
| groups | [Group](#gitlab-agent-gitlab-api-Group) | repeated |  |






<a name="gitlab-agent-gitlab-api-ProjectAccessCF"></a>

### ProjectAccessCF



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [int64](#int64) |  |  |
| roles | [string](#string) | repeated |  |






<a name="gitlab-agent-gitlab-api-User"></a>

### User



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [int64](#int64) |  |  |
| username | [string](#string) |  |  |





 

 

 

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

