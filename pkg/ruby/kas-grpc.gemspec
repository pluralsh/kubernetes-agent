Gem::Specification.new do |spec|
  spec.name          = 'kas-grpc'
  spec.version       = '0.1.0'
  spec.homepage      = 'https://gitlab.com/gitlab-org/cluster-integration/gitlab-agent'

  spec.summary       = 'Auto-generated gRPC client for KAS'
  spec.authors       = ['Tiger Watson', 'Hordur Freyr Yngvason', 'Timo Furrer']
  spec.email         = ['twatson@gitlab.com', 'hfyngvason@gitlab.com', 'tfurrer@gitlab.com']
  spec.license       = 'MIT'

  spec.files         = Dir['lib/**/*.rb']
  spec.require_paths = ['lib']

  spec.add_runtime_dependency 'grpc', '~> 1.0'
end
