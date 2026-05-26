FROM public.ecr.aws/sorah/ruby:4.0-dev-noble

RUN mkdir /app
COPY Gemfile Gemfile.lock /app/
WORKDIR /app
RUN bundle install
COPY . /app/
CMD ["bundle", "exec", "ruby", "server.rb"]
