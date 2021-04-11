## NAME
  gcp-hashicorp-packer-updater(1) -- update google source image versions in Packer templates

## SYNOPSIS

```
gcp-hashicorp-packer-updater 
   [-configuration name | -use-default-credentials] 
   [-project project] 
   [-filename path]
```

## DESCRIPTION
Updates the Google source image version in your packer template to the latest version.

If the packer template contains one or more builders of type `googlecompute` it will query
the Google Compute Engine API to determine the latest version of the specified source image.

If the `source_image_project_id` is not specified, it to derive this from the family or image name.
  If none is found the currently configured project or specified project will be used.

The utility assumes that the `source_image` follows the name pattern `.*-v.*`.

When a match is found, the utility will set or update the following fields:

- source_image
- source_image_family
- source_image_project_id

## EXAMPLES

```shell
$ gcp-hashicorp-packer-updater -filename tests/source-image.json 
2021/04/11 19:30:02 updating image from 'ubuntu-1804-bionic-v20210112' to 'ubuntu-1804-bionic-v20210325'
2021/04/11 19:30:02 setting image family to 'ubuntu-1804-lts'
2021/04/11 19:30:02 setting source image project to 'ubuntu-os-cloud'
```

## OPTIONS

* `-configuration name`
  the gcloud configuration to use for querying the Compute Engine API.

* `-use-default-credentials`
  use the Google default credentials from the environment.

* `-project project`
  to use, if not returned by the configuration or environment.
  
* `filename path`
  path of the packer template, default ./packer.json
  
## CAVEATS
- the utility only works with packer templates in JSON format.
- the utility assumes that the `source_image` has the name pattern `.*-v.*`.

## AUTHOR
Mark van Holsteijn

## COPYRIGHT
[binx.io B.V.](https://binx.io)

## SEE ALSO
- [How to keep your source image version in a Packer template up-to-date](https://binx.io/blog/2021/04/11/how-to-keep-your-source-image-version-in-a-packer-template-up-to-date/)