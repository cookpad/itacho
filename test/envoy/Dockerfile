FROM envoyproxy/envoy:v1.17.1

RUN apt-get update && \
  apt-get install -y curl gnupg2 &&\
  curl -sSfL https://sorah.jp/packaging/debian/3F0F56A8.pub.txt | gpg --dearmor > /usr/share/keyrings/sorah-rbpkg.gpg && \
  echo 'deb [signed-by=/usr/share/keyrings/sorah-rbpkg.gpg] https://cache.ruby-lang.org/lab/sorah/deb/ bionic main' > /etc/apt/sources.list.d/sorah-rbpkg.list && \
  apt-get update && \
  apt-get install -y ruby
COPY run.rb /run.rb

CMD ["ruby", "/run.rb"]
