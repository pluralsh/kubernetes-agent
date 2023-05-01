# Generated by the protocol buffer compiler.  DO NOT EDIT!
# Source: internal/module/notifications/rpc/rpc.proto for package 'gitlab.agent.notifications.rpc'

require 'grpc'
require 'internal/module/notifications/rpc/rpc_pb'

module Gitlab
  module Agent
    module Notifications
      module Rpc
        module Notifications
          class Service

            include ::GRPC::GenericService

            self.marshal_class_method = :encode
            self.unmarshal_class_method = :decode
            self.service_name = 'gitlab.agent.notifications.rpc.Notifications'

            rpc :GitPushEvent, ::Gitlab::Agent::Notifications::Rpc::GitPushEventRequest, ::Gitlab::Agent::Notifications::Rpc::GitPushEventResponse
          end

          Stub = Service.rpc_stub_class
        end
      end
    end
  end
end