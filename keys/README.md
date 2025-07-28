The keys are used for signing JWTs.

This service creates JWTs, so it needs the private key.
Any service that validates JWTs (i.e. any that uses JWTs) will need a public key.

We can put the public key inside the repo.
For ease of use during development, both keys are in `keys.rar`.
Please use those.

The keys themselves are mounted in `compose.yml`.
The script used to generate keys is `keygen.sh`.