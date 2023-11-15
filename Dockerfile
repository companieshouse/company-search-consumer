FROM 169942020521.dkr.ecr.eu-west-1.amazonaws.com/base/golang:1.19-bullseye-builder AS BUILDER

ARG SSH_PRIVATE_KEY
ARG SSH_PRIVATE_KEY_PASSPHRASE

RUN /lib/build

FROM 169942020521.dkr.ecr.eu-west-1.amazonaws.com/base/golang:debian11-runtime

COPY --from=BUILDER /build/out/app ./
