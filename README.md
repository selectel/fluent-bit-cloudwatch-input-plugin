# Fluent Bit input plugin for CloudWatch logs

## Quick start

```bash
docker run \
  --name fluent-bit-cloudwatch \
  --rm \
  -v ${PWD}/config/fluent-bit.yaml:/fluent-bit/etc/fluent-bit.yaml:ro \
  -v ${PWD}/sqlite:/var/lib/fluent-bit/cloudwatch/sqlite:rw \
  -e AWS_ACCESS_KEY_ID=your_access_key \
  -e AWS_SECRET_ACCESS_KEY=your_secret_key \
  ghcr.io/selectel/fluent-bit-cloudwatch-input-plugin:latest
```
