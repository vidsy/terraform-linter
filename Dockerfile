FROM vidsyhq/alpine
LABEL maintainer="Vidsy <tech@vidsy.co>"

ARG VERSION
LABEL version=$VERSION

ADD terraform-linter /usr/bin/terraform-linter
RUN chmod u+x /usr/bin/terraform-linter

CMD ["--tf-directory=."]
ENTRYPOINT ["terraform-linter"]
