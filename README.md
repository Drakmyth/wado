# Wado

A tool for mixing and munging Doom WAD files. It's a little eclectic and could go anywhere.

Where's Waldo?
Where's Wadlo?
Wado!

### Built With

* [![Golang][golang-shield]][golang-url]

## Instructions

### Development

The included [docker-compose configuration](./docker-compose.yaml) is set up to initialize a fully-featured dev environment utilizing [Air][air-url] for live-reloading and [Delve][delve-url] for remote debugging.

```
> docker compose up -d
```

To facilitate debugging in VSCode, a [launch configuration](./.vscode/launch.json) has been included that you can use to attach the VSCode debugger to the running dev environment.

### Release

To generate a release image, simply build the docker file. This utilizes a dedicated build stage along with the [scratch][scratch-url] Docker image to ensure the final image contains only the necessary resources and nothing else.

```
> docker build . -t go-template-release
```

## Example Application

The application provided with this example is a lightweight web service hosting static files out of the [./public](./public/) directory and templated responses out of the [./templates](./templates) directory. It is also set up to utilize environment variables for application configuration. These can either be set directly in the execution environment, included as part of a `docker-compose.yaml` definition, or provided in a `.env` file in the execution directory. The `.env` file will always be overridden by other methods so is particularly useful for development. The service will also fall back to a default value if one is not otherwise provided.

| Variable | Description                                                    | Default   |
| -------- | -------------------------------------------------------------- | --------- |
| `HOST`   | The host network that requests should be serviced from.        | localhost |
| `PORT`   | The port that hosts the configuration UI and aggregated feeds. | 80        |

### Example `.env` File

```
HOST=
PORT=8080
```

## License

Distributed under the MIT License. See [LICENSE.md](./LICENSE.md) for more information.


<!-- Reference Links -->
[golang-url]: https://go.dev
[golang-shield]: https://img.shields.io/badge/golang-09657c?style=for-the-badge&logo=go&logoColor=79d2fa
[air-url]: https://github.com/cosmtrek/air
[delve-url]: https://github.com/go-delve/delve
[scratch-url]: https://hub.docker.com/_/scratch/