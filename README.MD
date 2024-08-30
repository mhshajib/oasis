`warning`: **Development in Progress**: This project is currently under development. Expect changes and updates frequently.

<p align="center">
  <img src="oasis_logo.png" alt="Alt text" width="200"/>
</p>

<div align="center">
Oasis is a powerful CLI tool for generating Golang clean architecture modules including domain, repository, usecase, transformer, migration, seeder and config. Designed to streamline your development process, Oasis is an essential tool for Go developers.

<hr/>

![Go Version](https://img.shields.io/badge/go-1.x.x-blue)
![License](https://img.shields.io/badge/license-MIT-green)
![Build Status](https://img.shields.io/badge/build-passing-brightgreen)
[![All Contributors](https://img.shields.io/badge/all_contributors-0-orange.svg?style=flat-square)](#contributors-)

<hr/>

</div>

- [Installation](#installation)
- [Usage](#usage)
- [Features](#features)
- [Building Binaries](#building-binaries)
- [Contributing](#contributing)
- [Contributors](#contributors)
- [License](#license)

## Installation

#### Installation on Mac/Linux

```bash
curl https://raw.githubusercontent.com/mhshajib/oasis/main/install.sh -o install.sh \
&& chmod +x install.sh \
&& ./install.sh \
&& rm install.sh
```

#### Uninstallation

```bash
curl https://raw.githubusercontent.com/mhshajib/oasis/main/uninstall.sh -o uninstall.sh \
&& chmod +x uninstall.sh \
&& ./uninstall.sh \
&& rm uninstall.sh
```

#### Windows

**For windows download the binary and set environment variable so that you can access the binary from terminal**

**[Download Binary](https://github.com/mhshajib/oasis/releases)**

#### Usage

````markdown
## Usage

Quick start:

```bash
oasis generate module <module_name>
```
````

#### Features

```markdown
## Features

- **Module Generation**: Easily generate clean architecture modules like domain, repository, etc.
- **Customizable Templates**: Tailor the generated code to your needs.
- **Cross-Platform Support**: Works on various platforms.
```

## Building Binaries

Refer to the section below on [Automated Binary Builds](#automated-binary-builds) for details on building and accessing the latest binaries.

## Contributing

Contributions are welcome! Please read our [Contributing Guide](CONTRIBUTING.md) for more information.

## Contributors ✨

Thanks go to these wonderful people

<!-- ALL-CONTRIBUTORS-LIST:START -->
<!-- Do not remove or modify this section -->
<!-- ALL-CONTRIBUTORS-LIST:END -->

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Latest Releases

You can find the latest releases of Oasis, including pre-compiled binaries for various platforms, on our GitHub Releases page. Click here to access them: [Oasis Releases](https://github.com/mhshajib/oasis/releases)
