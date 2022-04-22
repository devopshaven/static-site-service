<img align="right" src="https://www.devopshaven.com/images/logo.svg" height="150">

# DevOps Haven - Static Site Service
Static site solver from MinIO bucket with single page application support for kubernetes environment.

---

[![License: LGPL v3](https://img.shields.io/badge/License-LGPL_v3-blue.svg)](https://www.gnu.org/licenses/lgpl-3.0)


**Container image:** `ghcr.io/devopshaven/static-site-service`

### Environment Variables

- `MINIO_ENDPOINT` - endpoint address of minio server without scheme eg: play.min.io
- `MINIO_ACCESS_KEY_ID` - minio access key id
- `MINIO_ACCESS_KEY_SECRET` - minio access key secret
- `MINIO_USE_SSL` - use ssl for the connection (true/false)
- `SITE_NAME` - the subdirectory name in the bucket where the site will be served from

The service currently supports only MinIO server.

<div style="text-align: right">Made with ❤️ at DevopsHaven Team</div>
