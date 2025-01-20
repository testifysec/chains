<!--
---
linkTitle: "Chains Configuration"
weight: 20
---
-->
# Chains Configuration

`Chains` works by observing `TaskRun` and `PipelineRun` executions, capturing relevant information, and storing it in a cryptographically-signed format.

`TaskRuns` and `PipelineRuns` can indicate inputs and outputs which are then captured and surfaced in the `Chains` payload formats where relevant.  
Chains uses the `Results` to _hint_ at the correct inputs and outputs. Check out [slsa-provenance.md](slsa-provenance.md) for more details.

## Chains Configuration

Chains uses a `ConfigMap` called `chains-config` in the `tekton-chains` namespace for configuration.  
Supported keys include:

### TaskRun Configuration

| Key                         | Description                                                                                                                                                                                      | Supported Values                                           | Default   |
| :-------------------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------- | :-------- |
| `artifacts.taskrun.format`  | The format to store `TaskRun` payloads in.                                                                                                                                                       | `in-toto`, `slsa/v1`, `slsa/v2alpha3`, `slsa/v2alpha4`      | `in-toto` |
| `artifacts.taskrun.storage` | The storage backend to store `TaskRun` signatures in. Multiple backends can be specified with a comma-separated list (e.g. `"tekton,oci"`). An empty string disables TaskRun artifacts. | `tekton`, `oci`, `gcs`, `docdb`, `grafeas`, `archivista` | `tekton`  |
| `artifacts.taskrun.signer`  | The signature backend to sign `TaskRun` payloads with.                                                                                                                                           | `x509`, `kms`                                              | `x509`    |

> **NOTE:**  
> - `slsa/v1` is an alias of `in-toto` for backwards compatibility.  
> - `slsa/v2alpha3` corresponds to the slsav1.0 spec and uses the latest [`v1` Tekton Objects](https://tekton.dev/docs/pipelines/pipeline-api/#tekton.dev/v1). Recommended for new Chains users who want the slsav1.0 spec.  
> - `slsa/v2alpha4` corresponds to the slsav1.0 spec and uses the latest [`v1` Tekton Objects](https://tekton.dev/docs/pipelines/pipeline-api/#tekton.dev/v1). It reads type-hinted results from [StepActions](https://tekton.dev/docs/pipelines/pipeline-api/#tekton.dev/v1alpha1.StepAction). Recommended for new Chains users who want the slsav1.0 spec.

### PipelineRun Configuration

| Key                                            | Description                                                                                                                                                                                                                                                                                 | Supported Values                                           | Default   |
| :--------------------------------------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | :--------------------------------------------------------- | :-------- |
| `artifacts.pipelinerun.format`                 | The format to store `PipelineRun` payloads in.                                                                                                                                                                                                                                              | `in-toto`, `slsa/v1`, `slsa/v2alpha3`, `slsa/v2alpha4`      | `in-toto` |
| `artifacts.pipelinerun.storage`                | The storage backend to store `PipelineRun` signatures in. Multiple backends can be specified with a comma-separated list (e.g. `"tekton,oci"`). An empty string disables PipelineRun artifacts.                                                                                    | `tekton`, `oci`, `gcs`, `docdb`, `grafeas`, `archivista` | `tekton`  |
| `artifacts.pipelinerun.signer`                 | The signature backend to sign `PipelineRun` payloads with.                                                                                                                                                                                                                                  | `x509`, `kms`                                              | `x509`    |
| `artifacts.pipelinerun.enable-deep-inspection` | This boolean option configures whether Chains should inspect child TaskRuns to capture inputs/outputs within a PipelineRun. `"false"` means only pipeline-level results are checked, whereas `"true"` means both pipeline and task level results are inspected.                           | `"true"`, `"false"`                                        | `"false"` |

> **NOTE:**  
> - For the Grafeas storage backend, currently only Container Analysis is supported. A configurable Grafeas server address is coming soon.  
> - `slsa/v1` is an alias of `in-toto` for backwards compatibility.  
> - `slsa/v2alpha3` corresponds to the slsav1.0 spec and uses the latest [`v1` Tekton Objects](https://tekton.dev/docs/pipelines/pipeline-api/#tekton.dev/v1). Recommended for new Chains users who want the slsav1.0 spec.  
> - `slsa/v2alpha4` corresponds to the slsav1.0 spec and uses the latest [`v1` Tekton Objects](https://tekton.dev/docs/pipelines/pipeline-api/#tekton.dev/v1). It reads type-hinted results from [StepActions](https://tekton.dev/docs/pipelines/pipeline-api/#tekton.dev/v1alpha1.StepAction) when `artifacts.pipelinerun.enable-deep-inspection` is set to `true`. Recommended for new Chains users who want the slsav1.0 spec.

### OCI Configuration

| Key                     | Description                                                                                                                                                                              | Supported Values                           | Default         |
| :---------------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | :----------------------------------------- | :-------------- |
| `artifacts.oci.format`  | The format to store `OCI` payloads in.                                                                                                                                                   | `simplesigning`                            | `simplesigning` |
| `artifacts.oci.storage` | The storage backend to store `OCI` signatures in. Multiple backends can be specified with a comma-separated list (e.g. `"oci,tekton"`). An empty string disables OCI artifacts.       | `tekton`, `oci`, `gcs`, `docdb`, `grafeas`   | `oci`           |
| `artifacts.oci.signer`  | The signature backend to sign `OCI` payloads with.                                                                                                                                       | `x509`, `kms`                              | `x509`          |

### KMS Configuration

| Key                  | Description                                                 | Supported Values                                                                                                                                | Default |
| :------------------- | :---------------------------------------------------------- | :---------------------------------------------------------------------------------------------------------------------------------------------- | :------ |
| `signers.kms.kmsref` | The URI reference to a KMS service for `KMS` signers.       | Supported schemes: `gcpkms://`, `awskms://`, `azurekms://`, `hashivault://`. See [Sigstore KMS Support](https://docs.sigstore.dev/cosign/kms_support) for details. |         |

### Storage Configuration

| Key                                  | Description                                                                                                                                                                                                                                                                                                         | Supported Values                                                                                                                                                                                                                                                                                                                                                                                                                                                    | Default  |
| :----------------------------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | :------- |
| `storage.gcs.bucket`                 | The GCS bucket for storage.                                                                                                                                                                                                                                                                                          |                                                                                                                                                                                                                                                                                                                                                                                                                                                                     |          |
| `storage.oci.repository`             | The OCI repository to store OCI signatures and attestations in.                                                                                                                                                                                                                                                     | If undefined and one of `artifacts.{oci,taskrun}.storage` includes OCI storage, attestations will be stored alongside the OCI artifact. Defining this value results in the OCI bundle stored in the designated location instead of alongside the image. See [cosign documentation](https://github.com/sigstore/cosign#specifying-registry) for details.                                                      |          |
| `storage.docdb.url`                  | The go-cloud URI reference to a docstore collection.                                                                                                                                                                                                                                                                 | `firestore://projects/[PROJECT]/databases/(default)/documents/[COLLECTION]?name_field=name`                                                                                                                                                                                                                                                                                                                                                                         |          |
| `storage.docdb.mongo-server-url` (optional) | The MongoDB connection URI, equivalent to the `MONGO_SERVER_URL` environment variable.                                                                                                                                                                                                                             | Example: `mongodb://[USER]:[PASSWORD]@[HOST]:[PORT]/[DATABASE]`                                                                                                                                                                                                                                                                                                                                                                                                     |          |
| `storage.docdb.mongo-server-url-dir` (optional) | The directory containing a file named `MONGO_SERVER_URL` with the MongoDB connection URI.                                                                                                                                                                                                                             | If the file `/mnt/mongo-creds-secret/MONGO_SERVER_URL` contains the MongoDB URL, set this to `/mnt/mongo-creds-secret`.                                                                                                                                                                                                                                                                                                  |          |
| `storage.docdb.mongo-server-url-path` (optional) | The file path that contains the MongoDB connection URI.                                                                                                                                                                                                                                                             | For example, if `/mnt/mongo-creds-secret/mongo-server-url` contains the MongoDB URL, set this to `/mnt/mongo-creds-secret/mongo-server-url`.                                                                                                                                                                                                                                                                             |          |
| `storage.grafeas.projectid`          | The project where the Grafeas server is located for storing occurrences.                                                                                                                                                                                                                                             |                                                                                                                                                                                                                                                                                                                                                                                                                                                                     |          |
| `storage.grafeas.noteid` (optional)  | The prefix for the note name used when creating a Grafeas note. Must be a string without spaces.                                                                                                                                                                                                                     |                                                                                                                                                                                                                                                                                                                                                                                                                                                                     |          |
| `storage.grafeas.notehint` (optional)| Sets the `human_readable_name` in the Grafeas ATTESTATION note. If not provided, defaults to `This attestation note was generated by Tekton Chains`.                                                                                                                                                              |                                                                                                                                                                                                                                                                                                                                                                                                                                                                     |          |
| `storage.archivista.url`         | The URL endpoint for the Archivista service.                                                                                                                                                                                                                                                                     | A valid HTTPS URL pointing to your Archivista instance (e.g. `https://archivista.testifysec.io`).                                                                                                                                                                                                                                                                                                                      | None |

#### docstore

For details on the go-cloud docstore URI format, see [Go Cloud Docstore](https://gocloud.dev/howto/docstore/). Chains supports the following docstore services:
- `firestore`
- `dynamodb`
- `mongo`

#### MongoDB

You can provide a MongoDB connection via multiple options:

- **Environment Variable:**  
  Set the connection URL in the `MONGO_SERVER_URL` environment variable in the Chains deployment.

- **ConfigMap Field (`storage.docdb.mongo-server-url`):**  
  Alternatively, set the connection URL in the `chains-config` ConfigMap. This field overrides the `MONGO_SERVER_URL` environment variable.

- **Directory Field (`storage.docdb.mongo-server-url-dir`):**  
  Set this to the directory containing a file named `MONGO_SERVER_URL`. This takes precedence over both the ConfigMap field and the environment variable.

- **File Path Field (`storage.docdb.mongo-server-url-path`):**  
  Directly reference the file containing the MongoDB URL. This field overrides all others.

**NOTE:**  
- When using the directory or file path fields, store the MongoDB URL in a secret and mount the secret so that Chains can pick up updates automatically.
- It is recommended to use `storage.docdb.mongo-server-url-dir` or `storage.docdb.mongo-server-url-path` rather than `storage.docdb.mongo-server-url` to avoid storing credentials in a ConfigMap.

#### Grafeas

For more information on Grafeas notes and occurrences, see [Grafeas Concepts](https://github.com/grafeas/grafeas/blob/master/docs/grafeas_concepts.md). To create occurrences, a note must first be created. Two types of occurrences are created:
- `ATTESTATION` Occurrence (note suffix: `-simplesigning`)
- `BUILD` Occurrence (note suffix: `-intoto`)  
If `storage.grafeas.noteid` is not set, the prefix `tekton-<NAMESPACE>` will be used.

### In-toto Configuration

| Key                         | Description                                    | Supported Values                                                                | Default                             |
| :-------------------------- | :--------------------------------------------- | :------------------------------------------------------------------------------ | :---------------------------------- |
| `builder.id`                | The builder ID for in-toto attestations.       |                                                                                 | `https://tekton.dev/chains/v2`      |
| `builddefinition.buildtype` | The buildType for in-toto attestations.        | `https://tekton.dev/chains/v2/slsa`, `https://tekton.dev/chains/v2/slsa-tekton` | `https://tekton.dev/chains/v2/slsa`   |

> **NOTE:**  
> - `builddefinition.buildtype` is valid for `slsa/v2alpha3` configurations only.  
> - Use `https://tekton.dev/chains/v2/slsa` for strict slsav1.0 compliance.  
> - Use `https://tekton.dev/chains/v2/slsa-tekton` for slsav1.0 with additional Tekton-specific details.

### Sigstore Features Configuration

#### Transparency Log

| Key                    | Description                                                        | Supported Values          | Default                      |
| :--------------------- | :----------------------------------------------------------------- | :------------------------ | :--------------------------- |
| `transparency.enabled` | Whether to enable automatic binary transparency uploads.         | `true`, `false`, `manual` | `false`                      |
| `transparency.url`     | The URL to upload binary transparency attestations to, if enabled. |                         | `https://rekor.sigstore.dev` |

**Note:** If `transparency.enabled` is set to `manual`, only TaskRuns and PipelineRuns with the annotation below will be uploaded to the transparency log:

```yaml
chains.tekton.dev/transparency-upload: "true"
```

#### Keyless Signing with Fulcio

| Key                                | Description                                                   | Supported Values                           | Default                                            |
| :--------------------------------- | :------------------------------------------------------------ | :----------------------------------------- | :------------------------------------------------- |
| `signers.x509.fulcio.enabled`      | Enable automatic certificates from Fulcio.                     | `true`, `false`                            | `false`                                            |
| `signers.x509.fulcio.address`      | Fulcio address for certificate requests.                      |                                           | `https://fulcio.sigstore.dev`                      |
| `signers.x509.fulcio.issuer`       | Expected OIDC issuer.                                          |                                           | `https://oauth2.sigstore.dev/auth`                 |
| `signers.x509.fulcio.provider`     | Provider for ID Token requests.                                | `google`, `spiffe`, `github`, `filesystem` | Unset (each provider will be attempted).           |
| `signers.x509.identity.token.file` | Path to file containing an ID Token.                           |                                           |                                                    |
| `signers.x509.tuf.mirror.url`      | TUF server URL; expects `$TUF_URL/root.json` to be present.      |                                           | `https://sigstore-tuf-root.storage.googleapis.com` |

#### KMS OIDC and Spire Configuration

| Key                               | Description                                                                                 | Supported Values | Default |
| :-------------------------------- | :------------------------------------------------------------------------------------------ | :--------------- | :------ |
| `signers.kms.auth.address`        | URI of the KMS server (e.g. `VAULT_ADDR`).                                                  |                  |         |
| `signers.kms.auth.token`          | Authentication token for the KMS server (e.g. `VAULT_TOKEN`).                               |                  |         |
| `signers.kms.auth.token-path`     | File path to store the KMS server Auth token (e.g. `/etc/kms-secrets`).                       |                  |         |
| `signers.kms.auth.oidc.path`      | Path used for OIDC authentication (e.g. `jwt` for Vault).                                   |                  |         |
| `signers.kms.auth.oidc.role`      | Role used for OIDC authentication.                                                         |                  |         |
| `signers.kms.auth.spire.sock`     | URI of the Spire socket for KMS token (e.g. `unix:///tmp/spire-agent/public/api.sock`).      |                  |         |
| `signers.kms.auth.spire.audience` | Audience for requesting a SVID from Spire.                                                  |                  |         |

> **NOTE:**  
> - If `signers.kms.auth.token-path` is set, create a secret and mount it to the specified path.  
> - To project updated token values without recreating pods, avoid using `subPath` in volume mounts.

### Visual Guide: ConfigMap Configuration Options

Refer to the diagram below to see a pictorial representation of signing and storage configuration options and how they relate to Chains artifacts.

![ConfigMap Configuration Diagram](../images/signing-storage-config-diagram.drawio.svg)

## Namespaces Restrictions in Chains Controller

Chains can be configured to monitor specific namespaces. If no namespaces are specified, the controller monitors all namespaces.

### Usage

Pass a comma-separated list of namespaces to the controller using the `--namespace` flag.

### Example

To restrict the controller to the `dev` and `test` namespaces, start the controller with:

```shell
--namespace=dev,test
```

In this example, the controller will only monitor TaskRuns and PipelineRuns in the `dev` and `test` namespaces.
```