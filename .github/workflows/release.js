// Quick and dirty script to use kontrol to package itself
const cp = require("child_process");
const fs = require("fs");
const path = require("path");

// Array of artifacts that look something like:
//  {
//      "name": "ghcr.io/frantjc/kontrol:latest",
//      "path": "ghcr.io/frantjc/kontrol:latest",
//      "goos": "linux",
//      "goarch": "amd64",
//      "goarm": "6",
//      "internal_type": 9,
//      "type": "Docker Image",
//      "extra": {
//          "DockerConfig": {
//              "goos": "linux",
//              "goarch": "amd64",
//              "goarm": "6",
//              "goamd64": "v1",
//              "dockerfile": "Dockerfile",
//              "image_templates": [
//                  "ghcr.io/frantjc/kontrol:{{ .Version }}",
//                  "ghcr.io/frantjc/kontrol:{{ .Major }}.{{ .Minor }}",
//                  "ghcr.io/frantjc/kontrol:{{ .Major }}",
//                  "ghcr.io/frantjc/kontrol:latest"
//              ],
//              "use": "docker"
//          }
//      }
//  }
const artifacts = JSON.parse(
    fs.readFileSync(
        path.join(__dirname, "../../dist/artifacts.json")
    )
);

// We only want the images as that is all kontrol operates on
artifacts.filter(artifact => artifact.type === "Docker Image").forEach(async artifact => {
    await new Promise(resolve => {
        // For each image built, package it with kontroller CRDs and roles
        const kontrol = cp.spawn(
            "kontrol",
            [
                "package",
                artifact.name,
                "--crds", path.join(__dirname, "../../manifests/frantj.cc_kontrollers.yaml"),
                "--roles", path.join(__dirname, "../../manifests/role.yaml")
            ],
            {
                stdio: 'inherit'
            }
        );

        kontrol.on('exit', (code) => {
            if (code > 0) {
                process.exit(code);
            }

            // Wherever this runs, it will be running just after the images were built,
            // so docker will be running, so the images will be cached in the docker
            // daemon, so we need to push them manually.
            const docker = cp.spawn(
                "docker",
                [
                    "push",
                    artifact.name
                ],
                {
                    stdio: 'inherit'
                }
            );

            docker.on('exit', (code) => {
                if (code > 0) {
                    process.exit(code);
                }

                resolve();
            });
        });
    });
});
