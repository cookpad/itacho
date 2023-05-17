# frozen_string_literal: true

require 'yaml'
require 'json'
require 'resolv'
require 'pp'
require 'erb'
$stdout.sync = true

statsd_exporter_ip = Resolv.getaddress('statsd-exporter')
p statsd_exporter_ip

config = YAML.load(ERB.new(File.read('/config.yaml.erb')).result)
File.open('/tmp/generated.json', 'w') { |f| f.puts(JSON.pretty_generate(config)) }

exec('envoy', '-c', '/tmp/generated.json', '--service-cluster', 'book', '--service-node', 'book', '--log-level', 'info')
