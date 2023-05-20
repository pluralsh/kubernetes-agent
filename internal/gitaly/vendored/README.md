Files in this directory were copied from Gitaly and modified to:

- Optimize memory consumption
- Handle empty repository case
- Reduce dependencies

Instruction on refreshing the vendored gRPC client and protobufs:

- Clone Gitaly repo somewhere. Use the directory name as `GITALY_DIR` env var value.
- Clone agent repo somewhere (this repo). Use the directory name as `AGENT_DIR` env var value.
- From the agent checkout directory run `GITALY_DIR=<dir> AGENT_DIR=<dir> ./build/vendor_gitaly.sh`. The script will copy the required protobuf and code files.
- Look at the diff and remove any bits of code that are not actually required. There are a few unrelated bits and pieces.
- Pay attention to dependencies that the code introduces, if any.
- Run `make regenerate-proto regenerate-mocks test lint` to see if everything works fine.
- Commit, push, open MR.
