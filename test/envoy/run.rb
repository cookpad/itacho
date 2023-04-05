# frozen_string_literal: true

require 'yaml'
require 'json'
require 'resolv'
require 'pp'
$stdout.sync = true

ip = Resolv.getaddress('statsd-exporter')
p ip

config = YAML.load_file('/config.yaml')
config['stats_sinks'][0]['config']['address']['socket_address']['address'] = ip
File.open('/tmp/generated.json', 'w') { |f| f.puts(JSON.pretty_generate(config)) }

exec('envoy', '-c', '/tmp/generated.json', '--bootstrap-version', '2', '--service-cluster', 'book', '--service-node', 'book',
'--log-level', 'info')
