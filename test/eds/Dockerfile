FROM public.ecr.aws/sorah/ruby:4.0-noble

WORKDIR /app
COPY Gemfile Gemfile.lock ./
RUN bundle install
COPY . ./
CMD ["bundle", "exec", "ruby", "eds.rb", "-p", "8080", "-o", "0.0.0.0"]
