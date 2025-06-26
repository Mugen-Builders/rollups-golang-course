<br>
<p align="center">
    <img src="https://github.com/user-attachments/assets/465b7615-842a-4f92-9f68-f3ffb8670fda" align="center" width="20%">
</p>
<br>
<div align="center">
    <i>A Linux-powered EVM rollup serving as a Debt Capital Market for the debtor economy</i>
</div>
<div align="center">
<b>Tokenized debt issuance through reverse campaign mechanism with collateralization</b>
</div>
<br>
<p align="center">
	<img src="https://img.shields.io/github/license/dcmhq/dcm?style=default&logo=opensourceinitiative&logoColor=white&color=959CD0" alt="license">
	<img src="https://img.shields.io/github/last-commit/dcmhq/dcm?style=default&logo=git&logoColor=white&color=D1DCCB" alt="last-commit">
</p>

##  Table of Contents

- [Overview](#overview)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Running](#running)
  - [Testing](#testing)

##  Overview

A debt capital market platform designed for the debtor economy, enabling content debtors to monetize their influence by issuing tokenized debt instruments collateralized. Through a reverse campaign mechanism, the platform connects debtors with investors, offering a structured and transparent alternative to finance scalable businesses while leveraging the economic potential of their audiences, ensuring legal compliance and attractive returns for investors.
<br>

[![Docs]][Link-docs] [![Deck]][Link-deck]
	
[Docs]: https://img.shields.io/badge/Documentation-959CD0?style=for-the-badge
[Link-docs]: https://docs.google.com/document/d/1l5D6sn9DBbaJFtTCfIM1gxoH7-10fVi9t2tsNr942Rw/edit?tab=t.0#heading=h.dfmi5re7vy34

[Deck]: https://img.shields.io/badge/Pitch%20Deck-D1DCCB?style=for-the-badge
[Link-deck]: https://www.canva.com/design/DAGVvlTnNpM/GsV9c1XuhYRYCrPK5811GA/view?utm_content=DAGVvlTnNpM&utm_campaign=designshare&utm_medium=link&utm_source=editor

##  Getting Started

###  Prerequisites
1. [Install Docker Desktop for your operating system](https://www.docker.com/products/docker-desktop/).

   To install Docker RISC-V support without using Docker Desktop, run the following command:
    
   ```shell
   docker run --privileged --rm tonistiigi/binfmt --install all
   ```

2. [Download and install the latest version of Node.js](https://nodejs.org/en/download).

3. Cartesi CLI is an easy-to-use tool to build and deploy your dApps on devnet. To install it, run:

   ```shell
   npm i -g @cartesi/cli
   ```

###  Running

1. Devnet mode:

   1.1 Build application:

   ```sh
   cartesi build
   ```

   1.2 Run application on devnet:

   ```sh
   cartesi run
   ```

2. Production mode:

- TODO

> [!NOTE]
> To reach the final state of the system for demo purposes, run the command below:
>
>   ```shell
>   make state
>   ```

###  Testing

```sh
make test
```
