import * as core from "@actions/core";
import * as tc from "@actions/tool-cache";
import * as cp from "@actions/exec";

import fs from "fs";
import path from "path";
import yaml from "yaml";

const packageJSON = JSON.parse(
  fs.readFileSync(path.join(__dirname, "../package.json")).toString()
);

const goreleaserYML = yaml.parse(
  fs.readFileSync(path.join(__dirname, "../../../.goreleaser.yaml")).toString()
);

async function run(): Promise<void> {
  try {
    const tool = "kontrol";
    const version = core.getInput("version") || packageJSON.version;

    // Turn RUNNER_ARCH into GOARCH.
    let arch;
    switch (process.env.RUNNER_ARCH) {
      case "X86":
      case "X64":
        arch = "amd64";
        break;
    }

    // Before we even attempt the download, check if goreleaser was configured
    // to build the GOARCH that we are trying to download.
    //
    // Note that this would become non-backwards-compatible if we remove support for
    // a GOARCH and it acts funny if we add support for one and someone uses it like so:
    //
    //  - uses: frantjc/kontrol@v0.1.0
    //    with:
    //      version: 0.2.0
    if (!goreleaserYML.builds[0].goarch.includes(arch)) {
      throw new Error(`unsupported architecture ${process.env.RUNNER_ARCH}`);
    }

    // Turn RUNNER_OS into GOOS.
    let os;
    switch (process.env.RUNNER_OS) {
      case "Linux":
        os = "linux";
        break;
      case "Windows":
        os = "windows";
        break;
      case "macOS":
        os = "darwin";
        break;
    }

    const versionOs = `${version}_${os}`;

    // Before we even attempt the download, check if goreleaser was configured
    // to build the GOOS that we are trying to download
    //
    // Note that this would become non-backwards-compatible if we remove support for
    // a GOOS and it acts funny if we add support for one and someone uses it like so:
    //
    //  - uses: frantjc/kontrol@v0.1.0
    if (!goreleaserYML.builds[0].goos.includes(os)) {
      throw new Error(`unsupported OS ${process.env.RUNNER_OS}`);
    }

    // Default to looking it up on PATH if install is explicitly set to false.
    let bin = tool;
    if (core.getBooleanInput("install")) {
      core.startGroup("install");

      // Look for kontrol in the cache.
      let dir = tc.find(tool, versionOs);

      // If we don't find kontrol in the cache, download, extract and cache it
      // from its GitHub release.
      if (!dir) {
        dir = await tc.cacheFile(
          path.join(
            await tc.extractTar(
              await tc.downloadTool(
                `https://github.com/frantjc/${tool}/releases/download/v${version}/${tool}_${version}_${os}_${arch}.tar.gz`
              )
            ),
            tool
          ),
          tool,
          tool,
          versionOs
        );
      }

      bin = path.join(dir, tool);

      core.addPath(dir);

      core.endGroup();
    }

    // Sanity check that kontrol was installed correctly.
    await cp.exec(bin, ["-v"]);
  } catch (err) {
    if (typeof err === "string" || err instanceof Error) core.setFailed(err);
  }
}

run();
