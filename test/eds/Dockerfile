FROM public.ecr.aws/sorah/ruby:3.2-jammy

WORKDIR /app
COPY Gemfile Gemfile.lock ./
RUN bundle install
COPY . ./
CMD ["bundle", "exec", "ruby", "eds.rb", "-p", "8080", "-o", "0.0.0.0"]
