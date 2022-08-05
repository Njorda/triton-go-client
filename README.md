# Triton Go client

[Triton](https://github.com/triton-inference-server) ships a generated(based upon the [protobuf](https://developers.google.com/protocol-buffers) definitions), which can be found [here](https://github.com/triton-inference-server/client/tree/main/src/grpc_generated/go). It should be possible to also [generate a http client](https://developers.google.com/protocol-buffers/docs/reference/go-generated) based upon the proto definitions. 

In this case we will utilize the model developed in the previous blog post [Trition with post and pre processing](https://github.com/Njorda/trition-ensemble). 

Since the go grpc client is not exported and instead have to be generated, I have copied it inside here to make it easier. When trying to generate the client the code seems to be out of date and not working any longer but made a [PR](https://github.com/triton-inference-server/client/pull/138) that is not yet reviewed with suggested changes. For now I will assume we have it generated locally. 


