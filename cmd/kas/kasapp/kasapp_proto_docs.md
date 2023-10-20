# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [cmd/kas/kasapp/kasapp.proto](#cmd_kas_kasapp_kasapp-proto)
    - [GatewayKasResponse](#gitlab-agent-kas-GatewayKasResponse)
    - [GatewayKasResponse.Error](#gitlab-agent-kas-GatewayKasResponse-Error)
    - [GatewayKasResponse.Header](#gitlab-agent-kas-GatewayKasResponse-Header)
    - [GatewayKasResponse.Header.MetaEntry](#gitlab-agent-kas-GatewayKasResponse-Header-MetaEntry)
    - [GatewayKasResponse.Message](#gitlab-agent-kas-GatewayKasResponse-Message)
    - [GatewayKasResponse.NoTunnel](#gitlab-agent-kas-GatewayKasResponse-NoTunnel)
    - [GatewayKasResponse.Trailer](#gitlab-agent-kas-GatewayKasResponse-Trailer)
    - [GatewayKasResponse.Trailer.MetaEntry](#gitlab-agent-kas-GatewayKasResponse-Trailer-MetaEntry)
    - [GatewayKasResponse.TunnelReady](#gitlab-agent-kas-GatewayKasResponse-TunnelReady)
    - [StartStreaming](#gitlab-agent-kas-StartStreaming)
  
- [Scalar Value Types](#scalar-value-types)



<a name="cmd_kas_kasapp_kasapp-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cmd/kas/kasapp/kasapp.proto



<a name="gitlab-agent-kas-GatewayKasResponse"></a>

### GatewayKasResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| tunnel_ready | [GatewayKasResponse.TunnelReady](#gitlab-agent-kas-GatewayKasResponse-TunnelReady) |  |  |
| header | [GatewayKasResponse.Header](#gitlab-agent-kas-GatewayKasResponse-Header) |  |  |
| message | [GatewayKasResponse.Message](#gitlab-agent-kas-GatewayKasResponse-Message) |  |  |
| trailer | [GatewayKasResponse.Trailer](#gitlab-agent-kas-GatewayKasResponse-Trailer) |  |  |
| error | [GatewayKasResponse.Error](#gitlab-agent-kas-GatewayKasResponse-Error) |  |  |
| no_tunnel | [GatewayKasResponse.NoTunnel](#gitlab-agent-kas-GatewayKasResponse-NoTunnel) |  |  |






<a name="gitlab-agent-kas-GatewayKasResponse-Error"></a>

### GatewayKasResponse.Error
Error represents a gRPC error that should be returned.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| status | [google.rpc.Status](#google-rpc-Status) |  | Error status as returned by gRPC. See https://cloud.google.com/apis/design/errors. |






<a name="gitlab-agent-kas-GatewayKasResponse-Header"></a>

### GatewayKasResponse.Header
Headers is a gRPC metadata.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| meta | [GatewayKasResponse.Header.MetaEntry](#gitlab-agent-kas-GatewayKasResponse-Header-MetaEntry) | repeated |  |






<a name="gitlab-agent-kas-GatewayKasResponse-Header-MetaEntry"></a>

### GatewayKasResponse.Header.MetaEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [gitlab.agent.prototool.Values](#gitlab-agent-prototool-Values) |  |  |






<a name="gitlab-agent-kas-GatewayKasResponse-Message"></a>

### GatewayKasResponse.Message
Message is a gRPC message data.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| data | [bytes](#bytes) |  |  |






<a name="gitlab-agent-kas-GatewayKasResponse-NoTunnel"></a>

### GatewayKasResponse.NoTunnel
No tunnels available at the moment.






<a name="gitlab-agent-kas-GatewayKasResponse-Trailer"></a>

### GatewayKasResponse.Trailer
Trailer is a gRPC trailer metadata.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| meta | [GatewayKasResponse.Trailer.MetaEntry](#gitlab-agent-kas-GatewayKasResponse-Trailer-MetaEntry) | repeated |  |






<a name="gitlab-agent-kas-GatewayKasResponse-Trailer-MetaEntry"></a>

### GatewayKasResponse.Trailer.MetaEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [gitlab.agent.prototool.Values](#gitlab-agent-prototool-Values) |  |  |






<a name="gitlab-agent-kas-GatewayKasResponse-TunnelReady"></a>

### GatewayKasResponse.TunnelReady
Tunnel is ready, can start forwarding stream.






<a name="gitlab-agent-kas-StartStreaming"></a>

### StartStreaming






 

 

 

 



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

