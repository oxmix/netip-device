# Netip Device Component
[![CI Status](https://github.com/oxmix/netip-device/workflows/Package%20release/badge.svg)](https://github.com/oxmix/netip-device/actions/workflows/package-release.yaml)

This repository is for public viewing and container assembly. For more information, follow the link https://cloudnetip.com/wiki

```shell
docker run -d --name netip.device --restart always \
  --cap-add=SYS_RAWIO --cap-add=SYS_ADMIN --uts=host \
  -e CONNECT_KEY=****** \
  --device /dev/md1:/dev/md1:r \
  -v /mnt/data/.netip-device:/_external/mnt/data/.netip-device:ro \
ghcr.io/oxmix/netip-device:latest
```
