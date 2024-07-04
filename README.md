# Wado

A tool for mixing and munging Doom WAD files. It's a little eclectic and could go anywhere.

### Built With

* [![Golang][golang-shield]][golang-url]

## Usage

Execute `wado help` for more detailed information.

Command  | Arguments                                      | Description
-------- | ---------------------------------------------- | -----------
analyze  | N/A                                            | *Not Yet Implemented*<br/>Analyze the difficulty of a WAD
convert  | `[flags] <input-wad-file> <output-wad-file>`   | Convert a WAD from Doom to Doom 2
generate | `[flags] <input-wad-folder> <output-wad-file>` | Generate a new WAD with random levels

## Development

### Build

```go
> go build
```

### Debugging

Using the [Delve][delve-url] debugger with CLI applications is a little tricky. See the [Delve documentation][delve-debug-url] for recommended procedures on how to do this.

### Release

While Wado is a CLI application and is intended to be executed by directly calling the binary, a containerized installation is also provided. This container utilizes a dedicated build stage along with the [scratch][scratch-url] Docker image to ensure the final image contains only the necessary resources and nothing else.

```
> docker build . -t wado
```

## License

Distributed under the MIT License. See [LICENSE.md](./LICENSE.md) for more information.


<!-- Reference Links -->
[golang-url]: https://go.dev
[golang-shield]: https://img.shields.io/badge/golang-09657c?style=for-the-badge&logo=go&logoColor=79d2fa
[delve-url]: https://github.com/go-delve/delve
[delve-debug-url]: https://github.com/go-delve/delve/blob/master/Documentation/faq.md#-how-can-i-use-delve-to-debug-a-cli-application
[scratch-url]: https://hub.docker.com/_/scratch/