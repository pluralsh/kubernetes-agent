# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [pkg/agentcfg/agentcfg.proto](#pkg_agentcfg_agentcfg-proto)
    - [AgentConfiguration](#gitlab-agent-agentcfg-AgentConfiguration)
    - [ChartCF](#gitlab-agent-agentcfg-ChartCF)
    - [ChartProjectSourceCF](#gitlab-agent-agentcfg-ChartProjectSourceCF)
    - [ChartSourceCF](#gitlab-agent-agentcfg-ChartSourceCF)
    - [ChartValuesCF](#gitlab-agent-agentcfg-ChartValuesCF)
    - [ChartValuesFileCF](#gitlab-agent-agentcfg-ChartValuesFileCF)
    - [CiAccessAsAgentCF](#gitlab-agent-agentcfg-CiAccessAsAgentCF)
    - [CiAccessAsCF](#gitlab-agent-agentcfg-CiAccessAsCF)
    - [CiAccessAsCiJobCF](#gitlab-agent-agentcfg-CiAccessAsCiJobCF)
    - [CiAccessAsImpersonateCF](#gitlab-agent-agentcfg-CiAccessAsImpersonateCF)
    - [CiAccessCF](#gitlab-agent-agentcfg-CiAccessCF)
    - [CiAccessGroupCF](#gitlab-agent-agentcfg-CiAccessGroupCF)
    - [CiAccessProjectCF](#gitlab-agent-agentcfg-CiAccessProjectCF)
    - [ConfigurationFile](#gitlab-agent-agentcfg-ConfigurationFile)
    - [ExtraKeyValCF](#gitlab-agent-agentcfg-ExtraKeyValCF)
    - [GitRefCF](#gitlab-agent-agentcfg-GitRefCF)
    - [GitopsCF](#gitlab-agent-agentcfg-GitopsCF)
    - [LoggingCF](#gitlab-agent-agentcfg-LoggingCF)
    - [ManifestProjectCF](#gitlab-agent-agentcfg-ManifestProjectCF)
    - [ObservabilityCF](#gitlab-agent-agentcfg-ObservabilityCF)
    - [PathCF](#gitlab-agent-agentcfg-PathCF)
    - [StarboardCF](#gitlab-agent-agentcfg-StarboardCF)
    - [StarboardFilter](#gitlab-agent-agentcfg-StarboardFilter)
    - [VulnerabilityReport](#gitlab-agent-agentcfg-VulnerabilityReport)
  
    - [log_level_enum](#gitlab-agent-agentcfg-log_level_enum)
  
- [Scalar Value Types](#scalar-value-types)



<a name="pkg_agentcfg_agentcfg-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## pkg/agentcfg/agentcfg.proto



<a name="gitlab-agent-agentcfg-AgentConfiguration"></a>

### AgentConfiguration
AgentConfiguration represents configuration for agentk.
Note that agentk configuration is not exactly the whole file as the file
may contain bits that are not relevant for the agent. For example, some
additional config for kas.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| gitops | [GitopsCF](#gitlab-agent-agentcfg-GitopsCF) |  |  |
| observability | [ObservabilityCF](#gitlab-agent-agentcfg-ObservabilityCF) |  |  |
| agent_id | [int64](#int64) |  | GitLab-wide unique id of the agent. |
| project_id | [int64](#int64) |  | Id of the configuration project. |
| ci_access | [CiAccessCF](#gitlab-agent-agentcfg-CiAccessCF) |  |  |
| starboard | [StarboardCF](#gitlab-agent-agentcfg-StarboardCF) |  |  |
| project_path | [string](#string) |  | Path of the configuration project |






<a name="gitlab-agent-agentcfg-ChartCF"></a>

### ChartCF



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| release_name | [string](#string) |  |  |
| source | [ChartSourceCF](#gitlab-agent-agentcfg-ChartSourceCF) |  |  |
| values | [ChartValuesCF](#gitlab-agent-agentcfg-ChartValuesCF) | repeated |  |
| namespace | [string](#string) | optional |  |
| max_history | [int32](#int32) | optional |  |






<a name="gitlab-agent-agentcfg-ChartProjectSourceCF"></a>

### ChartProjectSourceCF



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | Project id. e.g. gitlab-org/cluster-integration/gitlab-agent |
| path | [string](#string) |  | Path in the repository. e.g. charts/my-chart |
| ref | [GitRefCF](#gitlab-agent-agentcfg-GitRefCF) |  | Ref in the GitOps repository to fetch manifests from. |






<a name="gitlab-agent-agentcfg-ChartSourceCF"></a>

### ChartSourceCF



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project | [ChartProjectSourceCF](#gitlab-agent-agentcfg-ChartProjectSourceCF) |  |  |






<a name="gitlab-agent-agentcfg-ChartValuesCF"></a>

### ChartValuesCF



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| inline | [google.protobuf.Struct](#google-protobuf-Struct) |  |  |
| file | [ChartValuesFileCF](#gitlab-agent-agentcfg-ChartValuesFileCF) |  |  |






<a name="gitlab-agent-agentcfg-ChartValuesFileCF"></a>

### ChartValuesFileCF



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [string](#string) | optional | Project id. Can only be omitted if chart is coming from a GitLab project. In that case file is fetched from that project. e.g. gitlab-org/cluster-integration/gitlab-agent |
| ref | [GitRefCF](#gitlab-agent-agentcfg-GitRefCF) |  | Ref in the repository to fetch manifests from. |
| file | [string](#string) |  |  |






<a name="gitlab-agent-agentcfg-CiAccessAsAgentCF"></a>

### CiAccessAsAgentCF







<a name="gitlab-agent-agentcfg-CiAccessAsCF"></a>

### CiAccessAsCF



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| agent | [CiAccessAsAgentCF](#gitlab-agent-agentcfg-CiAccessAsAgentCF) |  |  |
| impersonate | [CiAccessAsImpersonateCF](#gitlab-agent-agentcfg-CiAccessAsImpersonateCF) |  |  |
| ci_job | [CiAccessAsCiJobCF](#gitlab-agent-agentcfg-CiAccessAsCiJobCF) |  | CiAccessAsCiUserCF ci_user = 4 [json_name = &#34;ci_user&#34;, (validate.rules).message.required = true]; |






<a name="gitlab-agent-agentcfg-CiAccessAsCiJobCF"></a>

### CiAccessAsCiJobCF







<a name="gitlab-agent-agentcfg-CiAccessAsImpersonateCF"></a>

### CiAccessAsImpersonateCF



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| username | [string](#string) |  |  |
| groups | [string](#string) | repeated |  |
| uid | [string](#string) |  |  |
| extra | [ExtraKeyValCF](#gitlab-agent-agentcfg-ExtraKeyValCF) | repeated |  |






<a name="gitlab-agent-agentcfg-CiAccessCF"></a>

### CiAccessCF
https://gitlab.com/gitlab-org/cluster-integration/gitlab-agent/-/blob/master/doc/kubernetes_ci_access.md


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| projects | [CiAccessProjectCF](#gitlab-agent-agentcfg-CiAccessProjectCF) | repeated |  |
| groups | [CiAccessGroupCF](#gitlab-agent-agentcfg-CiAccessGroupCF) | repeated |  |






<a name="gitlab-agent-agentcfg-CiAccessGroupCF"></a>

### CiAccessGroupCF



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| default_namespace | [string](#string) |  |  |
| access_as | [CiAccessAsCF](#gitlab-agent-agentcfg-CiAccessAsCF) |  |  |
| environments | [string](#string) | repeated |  |






<a name="gitlab-agent-agentcfg-CiAccessProjectCF"></a>

### CiAccessProjectCF



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| default_namespace | [string](#string) |  |  |
| access_as | [CiAccessAsCF](#gitlab-agent-agentcfg-CiAccessAsCF) |  |  |
| environments | [string](#string) | repeated |  |






<a name="gitlab-agent-agentcfg-ConfigurationFile"></a>

### ConfigurationFile
ConfigurationFile represents user-facing configuration file.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| gitops | [GitopsCF](#gitlab-agent-agentcfg-GitopsCF) |  |  |
| observability | [ObservabilityCF](#gitlab-agent-agentcfg-ObservabilityCF) |  | Configuration related to all things observability. This is about the agent itself, not any observability-related features. |
| ci_access | [CiAccessCF](#gitlab-agent-agentcfg-CiAccessCF) |  |  |
| starboard | [StarboardCF](#gitlab-agent-agentcfg-StarboardCF) |  |  |
| container_scanning | [StarboardCF](#gitlab-agent-agentcfg-StarboardCF) |  |  |






<a name="gitlab-agent-agentcfg-ExtraKeyValCF"></a>

### ExtraKeyValCF



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| val | [string](#string) | repeated | Empty elements are allowed by Kubernetes. |






<a name="gitlab-agent-agentcfg-GitRefCF"></a>

### GitRefCF
GitRef in the repository to fetch manifests from.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| tag | [string](#string) |  | A Git tag name, without `refs/tags/` |
| branch | [string](#string) |  | A Git branch name, without `refs/heads/` |
| commit | [string](#string) |  | A Git commit SHA |






<a name="gitlab-agent-agentcfg-GitopsCF"></a>

### GitopsCF



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| manifest_projects | [ManifestProjectCF](#gitlab-agent-agentcfg-ManifestProjectCF) | repeated |  |
| charts | [ChartCF](#gitlab-agent-agentcfg-ChartCF) | repeated |  |






<a name="gitlab-agent-agentcfg-LoggingCF"></a>

### LoggingCF



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| level | [log_level_enum](#gitlab-agent-agentcfg-log_level_enum) |  |  |
| grpc_level | [log_level_enum](#gitlab-agent-agentcfg-log_level_enum) | optional | optional to be able to tell when not set and use a different default value. |






<a name="gitlab-agent-agentcfg-ManifestProjectCF"></a>

### ManifestProjectCF
Project with Kubernetes object manifests.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) | optional | Project id. e.g. gitlab-org/cluster-integration/gitlab-agent |
| default_namespace | [string](#string) |  | Namespace to use if not set explicitly in object manifest. |
| paths | [PathCF](#gitlab-agent-agentcfg-PathCF) | repeated | A list of paths inside of the project to scan for .yaml/.yml/.json manifest files. |
| reconcile_timeout | [google.protobuf.Duration](#google-protobuf-Duration) |  | Reconcile timeout defines whether the applier should wait until all applied resources have been reconciled, and if so, how long to wait. |
| dry_run_strategy | [string](#string) |  | Dry run strategy defines whether changes should actually be performed, or if it is just talk and no action. https://github.com/kubernetes-sigs/cli-utils/blob/d6968048dcd80b1c7b55d9e4f31fc25f71c9b490/pkg/common/common.go#L68-L89 |
| prune | [bool](#bool) |  | Prune defines whether pruning of previously applied objects should happen after apply. |
| prune_timeout | [google.protobuf.Duration](#google-protobuf-Duration) |  | Prune timeout defines whether we should wait for all resources to be fully deleted after pruning, and if so, how long we should wait. |
| prune_propagation_policy | [string](#string) |  | Prune propagation policy defines the deletion propagation policy that should be used for pruning. https://github.com/kubernetes/apimachinery/blob/44113beed5d39f1b261a12ec398a356e02358307/pkg/apis/meta/v1/types.go#L456-L470 |
| inventory_policy | [string](#string) |  | InventoryPolicy defines if an inventory object can take over objects that belong to another inventory object or don&#39;t belong to any inventory object. This is done by determining if the apply/prune operation can go through for a resource based on the comparison the inventory-id value in the package and the owning-inventory annotation in the live object. https://github.com/kubernetes-sigs/cli-utils/blob/d6968048dcd80b1c7b55d9e4f31fc25f71c9b490/pkg/inventory/policy.go#L12-L66 |
| ref | [GitRefCF](#gitlab-agent-agentcfg-GitRefCF) |  | Ref in the GitOps repository to fetch manifests from. |






<a name="gitlab-agent-agentcfg-ObservabilityCF"></a>

### ObservabilityCF



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| logging | [LoggingCF](#gitlab-agent-agentcfg-LoggingCF) |  |  |






<a name="gitlab-agent-agentcfg-PathCF"></a>

### PathCF



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| glob | [string](#string) |  | Glob to use to scan for files in the repository. Directories with names starting with a dot are ignored. See https://github.com/bmatcuk/doublestar#about and https://pkg.go.dev/github.com/bmatcuk/doublestar/v2#Match for globbing rules. |






<a name="gitlab-agent-agentcfg-StarboardCF"></a>

### StarboardCF



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| vulnerability_report | [VulnerabilityReport](#gitlab-agent-agentcfg-VulnerabilityReport) |  |  |
| cadence | [string](#string) |  |  |






<a name="gitlab-agent-agentcfg-StarboardFilter"></a>

### StarboardFilter



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| namespaces | [string](#string) | repeated |  |
| resources | [string](#string) | repeated |  |
| containers | [string](#string) | repeated |  |
| kinds | [string](#string) | repeated |  |






<a name="gitlab-agent-agentcfg-VulnerabilityReport"></a>

### VulnerabilityReport



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| namespaces | [string](#string) | repeated |  |
| filters | [StarboardFilter](#gitlab-agent-agentcfg-StarboardFilter) | repeated |  |





 


<a name="gitlab-agent-agentcfg-log_level_enum"></a>

### log_level_enum


| Name | Number | Description |
| ---- | ------ | ----------- |
| info | 0 | default value must be 0 |
| debug | 1 |  |
| warn | 2 |  |
| error | 3 |  |


 

 

 



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

