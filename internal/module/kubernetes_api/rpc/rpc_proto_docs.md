# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [internal/module/kubernetes_api/rpc/rpc.proto](#internal_module_kubernetes_api_rpc_rpc-proto)
    - [ExtraKeyVal](#gitlab-agent-kubernetes_api-rpc-ExtraKeyVal)
    - [HeaderExtra](#gitlab-agent-kubernetes_api-rpc-HeaderExtra)
    - [ImpersonationConfig](#gitlab-agent-kubernetes_api-rpc-ImpersonationConfig)
  
    - [KubernetesApi](#gitlab-agent-kubernetes_api-rpc-KubernetesApi)
  
- [Scalar Value Types](#scalar-value-types)



<a name="internal_module_kubernetes_api_rpc_rpc-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## internal/module/kubernetes_api/rpc/rpc.proto



<a name="gitlab-agent-kubernetes_api-rpc-ExtraKeyVal"></a>

### ExtraKeyVal



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| val | [string](#string) | repeated |  |






<a name="gitlab-agent-kubernetes_api-rpc-HeaderExtra"></a>

### HeaderExtra
HeaderExtra is passed in grpctool.HttpRequest.extra.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| imp_config | [ImpersonationConfig](#gitlab-agent-kubernetes_api-rpc-ImpersonationConfig) |  |  |






<a name="gitlab-agent-kubernetes_api-rpc-ImpersonationConfig"></a>

### ImpersonationConfig
ImpersonationConfig is a representation of client-go rest.ImpersonationConfig.
See https://github.com/kubernetes/client-go/blob/release-1.22/rest/config.go#L201-L210


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| username | [string](#string) |  |  |
| groups | [string](#string) | repeated |  |
| uid | [string](#string) |  |  |
| extra | [ExtraKeyVal](#gitlab-agent-kubernetes_api-rpc-ExtraKeyVal) | repeated |  |





 

 

 


<a name="gitlab-agent-kubernetes_api-rpc-KubernetesApi"></a>

### KubernetesApi


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| MakeRequest | [.gitlab.agent.grpctool.HttpRequest](#gitlab-agent-grpctool-HttpRequest) stream | [.gitlab.agent.grpctool.HttpResponse](#gitlab-agent-grpctool-HttpResponse) stream | MakeRequest allows to make a HTTP request to Kubernetes API. |

 



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
