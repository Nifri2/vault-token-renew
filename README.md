# Vault Token Renewer

Renews vault service tokens as a Kubernetes CronJob

## How to Use

### Building

Download the latest [SOPS](https://github.com/getsops/sops) binary and add it to this Repos path as `./sops.bin`  

Then use `docker build -t yourregistry/yourtag:v1.0` 

### Using this Thing

You need to add the Tokens you want to renew to aa GitRepo as [SOPS encrypted YAML](https://github.com/getsops/sops#encrypting-using-hashicorp-vault) with the following Structure:
```yaml
tokens:
  - vault_url: https://vault.example.com
    token: s.1234567890abcdef
    name: svc-01
  - vault_url: https://vault.example.com
    token: s.1234567890abcdef
    name: svc-02
  - vault_url: https://vault.example.com
    token: s.1234567890abcdef
    name: svc-03
```

Then create a Kubernetes secret that contains the config for the programm (This can Onviosly be stored in vault aswell, aslong as it is availible for the programm):

```yaml
apiVersion: v1
data:
  GIT_TOKEN: yourgit project token
  GIT_URL: URL to the Git Project where the encrypted tokens are stored
  VAULT_TOKEN: vault token to update the tokens
  VAULT_TRANSIT_KEY: Key used to decrypt the yaml (like /v1/transit/keys/yourtoken)
  VAULT_URL: ULR of the vault where the transit Key is stred
kind: Secret
metadata:
  name: renew-secret
  namespace: somens
```

Then you can simply create a cronjob that mounts the secret as Envirnment variables

```yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: token-renewer-job
  namespace: hc-vault
spec:
  concurrencyPolicy: Allow
  failedJobsHistoryLimit: 1
  jobTemplate:
    spec:
      template:
        spec:
          containers:
            - envFrom:
                - secretRef:
                    name: renew-secret
                    optional: false
              image: tokenrenewer:v1.0.1
              imagePullPolicy: Always
              name: renewer
          dnsPolicy: ClusterFirst
          restartPolicy: Never
  schedule: 0 * * * *
  successfulJobsHistoryLimit: 3
```