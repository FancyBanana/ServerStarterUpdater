# What is ServerStarterUpdater?

It is a small utility I wrote to automatically update `server-setup-config.yaml`, which is used by ServerStarter, 
to the latest CurseForge modpack release. It is designed to be run just before launching ServerStarter.

**ServerStarterUpdater** works by finding the download link for the latest release of a modpack on CurseForge and updates the `install.modpackUrl` entry with the new value.
It supports mod search by mod ID *or* by mod slug, you can use only one. 
