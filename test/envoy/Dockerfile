FROM envoyproxy/envoy:v1.9.0

RUN apt-get update && apt-get install -y software-properties-common curl
RUN apt-add-repository ppa:brightbox/ruby-ng
RUN apt-get update && apt-get install -y ruby2.4
COPY run.rb /run.rb

CMD ["ruby", "/run.rb"]
