FROM public.ecr.aws/sorah/ruby:3.2-dev-jammy

RUN mkdir /app
COPY Gemfile Gemfile.lock /app/
WORKDIR /app
RUN bundle install
COPY . /app/
CMD ["bundle", "exec", "ruby", "server.rb"]
