require 'resolv'
require 'sinatra'

get '/' do
  sleep ENV['SLEEP'].to_f
  raise 'error' if ENV['ERROR_RATE'] && rand(0..ENV['ERROR_RATE'].to_i).zero?

  response = "GET,#{env['HTTP_HOST']},#{ENV['RESPONSE']}"
  response << ",x-test-header=#{request.get_header('HTTP_X_TEST')}" if request.has_header?('HTTP_X_TEST')
  response
end

post '/' do
  raise 'error' if rand(0..ENV['ERROR_RATE'].to_i).zero?

  "POST and #{ENV['RESPONSE'] || 'hello'}"
end

get '/ip/:name' do
  Resolv.getaddress(params[:name])
end
