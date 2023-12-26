# `kubepf`: Handy port-forwarding for k8s


![MIT License](https://img.shields.io/badge/license-MIT-_red.svg)
![Go Report](https://goreportcard.com/badge/github.com/alpkeskin/kubepf)
![Release](https://img.shields.io/github/release/alpkeskin/kubepf)

**kubepf** is a simple utility for creating and managing port-forwarding in k8s. It's written in Go and uses Cobra for CLI.


Here is a `kubepf` demo:
![kubepf](assets/kubepf.gif)

## Installation

```bash
go install -v github.com/alpkeskin/kubepf/cmd/kubepf@latest
```

## Configuration
`kubepf` uses a config file named `.kubepf` in your **home directory**. You can create it manually or use `kubepf init` command to create it. Here is an example config file:

```yaml
# .kubepf config file. Edit it.
projects:
  - name: project1
    namespace: namespace1
    services:
      - name: service1
        local_port: 8081
        target_port: 8081
      - name: service2
        local_port: 8082
        target_port: 8082

  - name: project2
    namespace: namespace2
    services:
      - name: service3
        local_port: 8083
        target_port: 8083
      - name: service4
        local_port: 8084
        target_port: 8084
```

## Usage
**List projects and services in .kubepf config file**
```bash
kubepf list
```

**Start port-forwarding for project**
```bash
kubepf <project_name>
```

**List active port-forwarding**
```bash
kubepf active
```

**Kill port-forwarding for project**
```bash
kubepf kill <project_name>
```

**Kill port-forwarding for service**
```bash
kubepf kill <service_name>
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.