FROM public.ecr.aws/lambda/provided:al2 as build
RUN yum install -y golang
RUN go env -w GOPROXY=direct
ADD go.mod go.sum ./
RUN go mod download
ADD . .
RUN go build -o /main

FROM public.ecr.aws/lambda/provided:al2

LABEL maintainer="rluisr" \
  org.opencontainers.image.created=$BUILD_DATE \
  org.opencontainers.image.url="https://github.com/rluisr/ecssc" \
  org.opencontainers.image.source="https://github.com/rluisr/ecssc"
  org.opencontainers.image.version=$VERSION \
  org.opencontainers.image.revision=$VCS_REF \
  org.opencontainers.image.vendor="rluisr" \
  org.opencontainers.image.title="ecssc" \
  org.opencontainers.image.description="ecssc(ECS State Check) is a Lambda function for notification to Slack if the ECS task event is changed." \
  org.opencontainers.image.licenses="WTFPL"

COPY --from=build /main /main
ENTRYPOINT [ "/main" ]   
