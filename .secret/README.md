# secrets

## usage

1. build secret.json
2. upload secret.json to secret-manager

```bash
cd development # path/to/secret-files
node ../build.js # output: secret.json
gcloud secrets versions add ${SECRET_NAME} --project=${PROJECT} --data-file=secret.json
```

### first time

create new secret

```bash
gcloud secrets create ${SECRET_NAME} --project=${PROJECT} --replication-policy=automatic --data-file=secret.json
```

## jwt keys

see: [gist](https://gist.github.com/maxogden/62b7119909a93204c747633308a4d769)

```bash
# ES512
# private key
openssl ecparam -genkey -name secp521r1 -noout -out ecdsa-p521-private.pem

# public key
openssl ec -in ecdsa-p521-private.pem -pubout -out ecdsa-p521-public.pem
```
