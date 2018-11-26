# frozen_string_literal: true

require 'cgi'
require 'net/http'
require 'uri'
require 'json'
require 'rack'

lib = File.expand_path('grpc/lib', __dir__)
$LOAD_PATH.unshift(lib) unless $LOAD_PATH.include?(lib)

require 'grpc'
require 'health_services_pb'

def raise_error(response = nil)
  p response if response
  p response.body if response
  raise('invalid response')
end

envoy_url = URI('http://localhost:9211')
app_url = URI('http://localhost:3081')
sds_url = URI('http://localhost:4000')
ab_testing_app = nil

catch(:break) do
  i = 0
  loop do
    begin
      Net::HTTP.start(app_url.host, app_url.port) do |http|
        response = http.get('/ip/ab-testing-app')
        ab_testing_app = response.body
        throw(:break)
      end
    rescue EOFError, SystemCallError
      raise('Can not run the app container') if i == 19 # Overall retries end within 3.8s.

      puts 'waiting the app container to run...'
      sleep((2 * i) / 100.0)
      i += 1
    end
  end
end

catch(:break) do
  i = 0
  loop do
    begin
      Net::HTTP.start(sds_url.host, sds_url.port) do |http|
        http.get('/v1/registration/dummy')
        throw(:break)
      end
    rescue EOFError, SystemCallError
      raise('Can not run the app container') if i == 19 # Overall retries end within 3.8s.

      puts 'waiting the sds container to run...'
      sleep((2 * i) / 100.0)
      i += 1
    end
  end
end

catch(:break) do
  i = 0
  loop do
    begin
      Net::HTTP.start(envoy_url.host, envoy_url.port) do |http|
        response = http.get('/')
        throw(:break) if response.code == '404'
      end
    rescue EOFError, SystemCallError
      raise('Can not run the envoy container') if i == 19 # Overall retries end within 3.8s.

      puts 'waiting the envoy container to run...'
      sleep((2 * i) / 100.0)
      i += 1
    end
  end
end

puts 'register hosts to legacy sds'
Net::HTTP.start(sds_url.host, sds_url.port) do |http|
  payload = { ip: ab_testing_app, port: 8080, revision: 'a', tags: { az: 'b', region: 'ap-northeast-1', instance_id: 'test-instance' }.to_json }
  response = http.post('/v1/registration/ab-testing-development', payload.map { |k, v| "#{k}=#{CGI.escape(v.to_s)}" }.join('&'))
  puts response.code, response.body
  raise_error if response.code != '200'

  response = http.get('/v1/registration/ab-testing-development')
  puts response.code, response.body
  raise_error if response.code != '200'
  raise_error if JSON.parse(response.body)['hosts'].size != 1

  puts 'pass'
end

puts 'ensure Envoy has healty hosts'
Net::HTTP.start(envoy_url.host, envoy_url.port) do |http|
  catch(:break) do
    i = 0
    loop do
      stub = Grpc::Health::V1::Health::Stub.new(
        "#{envoy_url.host}:#{envoy_url.port}",
        :this_channel_is_insecure,
        channel_args: { 'grpc.default_authority' => 'ab-testing' }
      )
      begin
        response = stub.check(Grpc::Health::V1::HealthCheckRequest.new(service: 'test'))
        response_user = http.get('/', 'Host' => 'user')
        throw(:break) if response && response.status == :SERVING && response_user.code == '200'
      rescue GRPC::Cancelled, GRPC::Unavailable => e
        p e
      end

      raise('Can not fetch healty upstreams') if i > 30

      puts 'waiting the envoy to fetch from SDS...'
      sleep((2 * i) / 100.0)
      i += 1
    end
  end

  sleep 1

  response = http.get('/', 'Host' => 'user')
  raise_error if response.code != '200'
  raise_error if response.body != 'GET,user-service,user'

  stub = Grpc::Health::V1::Health::Stub.new(
    "#{envoy_url.host}:#{envoy_url.port}",
    :this_channel_is_insecure,
    channel_args: { 'grpc.default_authority' => 'ab-testing' }
  )
  response = stub.check(Grpc::Health::V1::HealthCheckRequest.new(service: 'test'))
  raise_error unless response.status == :SERVING

  puts 'pass'
end

puts 'OK'
