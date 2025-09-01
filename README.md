# 🦥 slothctl 🌳

Your friendly infrastructure automation assistant!

`slothctl` is a powerful command-line interface (CLI) tool, meticulously crafted in Go, to simplify and automate complex infrastructure management tasks. Like a sloth, it moves slowly but deliberately, ensuring everything is done correctly and without rush, handling the heavy lifting for you.

Our little sloth climbs the trees of your infrastructure, managing everything from virtual containers and configuration to secrets and IT support tickets.

## 📜 What is this?

`slothctl` is the central nervous system for an opinionated infrastructure stack. It provides a unified interface to orchestrate and interact with various services:

-   **Incus:** Manages system containers and VMs for lightweight virtualization.
-   **SaltStack:** Applies configurations and states to your servers, ensuring consistency.
-   **HashiCorp Vault:** Secures, stores, and tightly controls access to tokens, passwords, certificates, and other secrets.
-   **GLPI:** Integrates with your IT asset management and ticketing system, allowing you to manage support tickets directly from the command line.

## ✨ Core Features

-   **🚀 Bootstrap Infrastructure:** Initialize a complete environment from scratch.
-   **🖥️ Server Management:** Register, configure, list, and connect to servers via SSH.
-   **🔒 VPN Control:** Easily manage VPN connections and configurations.
-   **🎫 GLPI Ticketing:** Create, update, and manage GLPI tickets without leaving your terminal.
-   **🤫 Secrets Management:** Seamless integration with Vault for secure operations.
-   **⚙️ Centralized Configuration:** A single source of truth for managing the entire stack.

## 📦 Installation

You can install `slothctl` using the provided shell script, which handles all dependencies and setup:

```bash
./install_slothctl.sh
```

Alternatively, if you have a Go environment set up, you can build it from the source.

## 🚀 Getting Started

After installation, most configuration is handled by the tool itself. The primary configuration file is typically located at `~/.config/slothctl/config.yaml`.

## Examples

`slothctl` provides a rich set of commands to interact with your infrastructure.

### Managing Servers

List all registered servers:

```bash
slothctl server list
```

SSH into a managed server:

```bash
slothctl server ssh <server-name>
```

Execute a command on a server without a full SSH session:

```bash
slothctl server ssh <server-name> --exec "uptime"
```

### Managing VPN

Connect to the default VPN profile:

```bash
slothctl vpn connect
```

Check the current VPN status:

```bash
slothctl vpn status
```

Disconnect from the VPN:

```bash
slothctl vpn disconnect
```

### Interacting with GLPI Tickets

Create a new ticket:

```bash
slothctl glpi tickets create --title "Server is slow" --description "The main web server is experiencing high latency." --urgency 5
```

List your open tickets:

```bash
slothctl glpi tickets list
```

---

> The sloth sees the path and follows it. Let `slothctl` be your guide.
