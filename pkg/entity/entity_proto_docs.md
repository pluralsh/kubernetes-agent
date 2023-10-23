# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [pkg/entity/entity.proto](#pkg_entity_entity-proto)
    - [AgentMeta](#gitlab-agent-entity-AgentMeta)
    - [GitalyRepository](#gitlab-agent-entity-GitalyRepository)
    - [KubernetesVersion](#gitlab-agent-entity-KubernetesVersion)
  
- [Scalar Value Types](#scalar-value-types)



<a name="pkg_entity_entity-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## pkg/entity/entity.proto



<a name="gitlab-agent-entity-AgentMeta"></a>

### AgentMeta
AgentMeta contains information about agentk.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| version | [string](#string) |  | Version of the binary. |
| commit_id | [string](#string) |  | Short commit sha of the binary. |
| pod_namespace | [string](#string) |  | Namespace of the Pod running the binary. |
| pod_name | [string](#string) |  | Name of the Pod running the binary. |
| kubernetes_version | [KubernetesVersion](#gitlab-agent-entity-KubernetesVersion) |  | Version of the Kubernetes cluster. |






<a name="gitlab-agent-entity-GitalyRepository"></a>

### GitalyRepository
Modified copy of Gitaly&#39;s https://gitlab.com/gitlab-org/gitaly/-/blob/55cb537898bce04e5e44be074a4d3d441e1f62b6/proto/shared.proto#L25
Was copied to avoid exposing Gitaly type in the API and forcing the consumer to have a dependency on Gitaly.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| storage_name | [string](#string) |  |  |
| relative_path | [string](#string) |  |  |
| git_object_directory | [string](#string) |  | Sets the GIT_OBJECT_DIRECTORY envvar on git commands to the value of this field. It influences the object storage directory the SHA1 directories are created underneath. |
| git_alternate_object_directories | [string](#string) | repeated | Sets the GIT_ALTERNATE_OBJECT_DIRECTORIES envvar on git commands to the values of this field. It influences the list of Git object directories which can be used to search for Git objects. |
| gl_repository | [string](#string) |  | Used in callbacks to GitLab so that it knows what repository the event is associated with. May be left empty on RPC&#39;s that do not perform callbacks. During project creation, `gl_repository` may not be known. |
| gl_project_path | [string](#string) |  | The human-readable GitLab project path (e.g. gitlab-org/gitlab-ce). When hashed storage is use, this associates a project path with its path on disk. The name can change over time (e.g. when a project is renamed). This is primarily used for logging/debugging at the moment. |






<a name="gitlab-agent-entity-KubernetesVersion"></a>

### KubernetesVersion
Version information of the Kubernetes cluster.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| major | [string](#string) |  | Major version of the Kubernetes cluster. |
| minor | [string](#string) |  | Minor version of the Kubernetes cluster. |
| git_version | [string](#string) |  | Git version of the Kubernetes cluster. |
| platform | [string](#string) |  | Platform of the Kubernetes cluster. |





 

 

 

 



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

