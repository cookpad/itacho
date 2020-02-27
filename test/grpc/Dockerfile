FROM ruby:2.6

RUN mkdir /app
COPY Gemfile Gemfile.lock /app/
WORKDIR /app
RUN bundle install
COPY . /app/
CMD ["bundle", "exec", "ruby", "server.rb"]
