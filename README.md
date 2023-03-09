# JSON Patcher

## Description 

`json-patcher` is a CLI tools to applying a list of JSON patch (rfc6902) from a declarative config file. Think of it as [kustomize patch](https://kubectl.docs.kubernetes.io/references/kustomize/builtins/#_patchesjson6902_) but for arbitrary JSON files. 

At it's core it is using the ["github.com/evanphx/json-patch/v5"](https://github.com/evanphx/json-patch) library.

## Install 

Download the binaries from the [releases](https://github.com/vdbulcke/json-patcher/releases) page. 

### Validate Signature With Cosign

Make sure you have `cosign` installed locally (see [Cosign Install](https://docs.sigstore.dev/cosign/installation/)).


Then you can use the `./verify_signature.sh` in this repo: 

```bash
./verify_signature.sh PATH_TO_DOWNLOADED_ARCHIVE TAG_VERSION
```
for example
```bash
$ ./verify_signature.sh ~/Downloads/json-patcher_0.1.0_Linux_x86_64.tar.gz v0.1.0

Checking Signature for version: v0.1.0
Verified OK

```


### Add To PATH 



```bash
sudo mv json-patcher /usr/local/bin/
```

## Getting Started



Create a patch file `patch.yaml`: 
```yaml
---
## List of patches
patches:

## Patch on a source and a destination
- 
  ## source:  same as '{}' as source file
  source: NEW
  ## destination: where the json should be written after 
  ##              all patches have been applied
  destination: ./generated.json
  json_patch: |-
    ## this is a first patch 
    - op: add
      path: "/foo"
      value: "baz"
    ## this is a second patch 
    - op: add
      path: "/hello"
      value: "world"

## You can add here another list of json_patch to another sources and/or destinations
```


Apply the patch
```bash
json-patcher apply -p patch.yaml
```

See the result
```bash
$ cat generated.json 
{"foo":"baz","hello":"world"}
```

See [./example/patch.yaml](./example/patch.yaml) for details information about configuration of patches.


## CLI Usage 

See [json-patcher](./doc/json-patcher.md) for CLI usage.

## Completion 

See [json-patcher completion](./doc/json-patcher_completion.md).

