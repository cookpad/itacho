FROM ruby:2.5

RUN mkdir /app
COPY Gemfile Gemfile.lock /app/
WORKDIR /app
RUN bundle install
COPY app.rb /app/
CMD ["bundle", "exec", "ruby", "app.rb", "-p", "8080", "-o", "0.0.0.0"]
