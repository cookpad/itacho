require 'sinatra'
require 'sinatra/json'
require 'rack/contrib/post_body_content_type_parser'
require 'securerandom'

REGISTRY = {}
Entry = Data.define(:ip, :port)

# Mock EDS API
use Rack::PostBodyContentTypeParser
post '/:version/discovery\:endpoints' do
  type = {
    'v2' => 'type.googleapis.com/envoy.api.v2.ClusterLoadAssignment',
    'v3' => 'type.googleapis.com/envoy.config.endpoint.v3.ClusterLoadAssignment',
  }.fetch(params[:version])
  resources = (params[:resource_names] || []).map do |resource_name|
    {
      '@type': type,
      cluster_name: resource_name,
      endpoints: [{
        lb_endpoints: REGISTRY.fetch(resource_name, []).map { |entry|
          {
            endpoint: {
              address: {
                socket_address: {
                  address: entry.ip,
                  port_value: entry.port,
                },
              },
            },
          }
        },
      }],
    }
  end

  json({
    version_info: SecureRandom.uuid,
    resources:,
  })
end

# Used from test.rb
post '/v1/registration/:name' do
  name = params[:name]
  REGISTRY[name] ||= []
  REGISTRY[name] << Entry.new(
    ip: params[:ip],
    port: params[:port],
  )

  200
end

get '/v1/registration/:name' do
  entries = REGISTRY.fetch(params[:name], [])
  json({
    hosts: entries.map { |e| { ip_address: e.ip, port: e.port } }
  })
end
