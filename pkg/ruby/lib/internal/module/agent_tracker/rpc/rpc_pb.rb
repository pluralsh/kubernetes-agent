# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: internal/module/agent_tracker/rpc/rpc.proto

require 'google/protobuf'

require 'internal/module/agent_tracker/agent_tracker_pb'
Google::Protobuf::DescriptorPool.generated_pool.build do
  add_file("internal/module/agent_tracker/rpc/rpc.proto", :syntax => :proto3) do
    add_message "gitlab.agent.agent_tracker.rpc.GetConnectedAgentsRequest" do
      oneof :request do
        optional :project_id, :int64, 1, json_name: "projectId"
        optional :agent_id, :int64, 2, json_name: "agentId"
      end
    end
    add_message "gitlab.agent.agent_tracker.rpc.GetConnectedAgentsResponse" do
      repeated :agents, :message, 1, "gitlab.agent.agent_tracker.ConnectedAgentInfo", json_name: "agents"
    end
  end
end

module Gitlab
  module Agent
    module AgentTracker
      module Rpc
        GetConnectedAgentsRequest = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("gitlab.agent.agent_tracker.rpc.GetConnectedAgentsRequest").msgclass
        GetConnectedAgentsResponse = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("gitlab.agent.agent_tracker.rpc.GetConnectedAgentsResponse").msgclass
      end
    end
  end
end