# Go-Pub
A server-side implementation of the [ActivityPub](https://www.w3.org/TR/activitypub/) social networking protocol.

Special thanks to [tedu](https://www.tedunangst.com/) for the knowledge and inspiration.

## Useful Resources
- ActivityPub
    - [ActivityPub Specification](https://www.w3.org/TR/activitypub/)
    - [ActivityPub Vocabulary](https://www.w3.org/TR/activitystreams-vocabulary/)
- JSON-LD
    - [JSON-LD](https://json-ld.org/)
    - [Go Module](https://github.com/cheebz/arb) created for handling arbitrary JSON
- HTTP Signatures
    - [Spec](https://datatracker.ietf.org/doc/html/draft-cavage-http-signatures)
    - [Go Module](https://github.com/cheebz/sigs) created for signing and validating requests

## Configuration
See **.env.example** for an example configuration.

### Notes
- AUTH - Authorization endpoint. GET request will be made to this endpoint to authorize requests, as necessary.
- CLIENT - Requests made without the "application/activity+json" Accept header will be reverse proxied to this URL. Can also provide a directory path here to serve static files.
- RSA_PUBLIC_KEY/RSA_PRIVATE_KEY - Paths to RSA public and private keys, respectively. Used to sign requests for federation.

*Currently the application supports only PostgreSQL databases (hoping to add more eventually). Execute the init_db.sql statement to build the required tables.*

## Docker
Container image can be created by running `docker build -t cheebz/go-pub .` or `docker build -t cheebz/go-pub -f Dockerfile.prod .`

Docker-compose files are available that include PostgreSQL and Redis images.

## Kubernetes
There is a Helm chart in the [go-pub-k8s](https://github.com/cbodonnell/go-pub-k8s) repository.
