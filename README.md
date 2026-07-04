## Boids implementation written in Go + Sysinfo fetching

- Fetches system information.
- This uses the tcell library for terminal rendering.
- The mini project was developed for learning purposes and to explore the Go programming language.

### Installation
Requires [Go](https://go.dev/dl/) 1.25 or later.

```bash
go install github.com/aeron-cap/boidfetch@latest
```

This builds the binary and places it in `$GOBIN` (or `$GOPATH/bin` if `$GOBIN` isn't set). Make sure that directory is on your `PATH`:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

Add that line to your `~/.bashrc` (or `~/.zshrc`) to make it persistent, then reload:

```bash
source ~/.bashrc
```

Verify it's installed:

```bash
boidfetch
```

### Building from source

```bash
git clone https://github.com/aeron-cap/boidfetch.git
cd boidfetch
go build -o boidfetch .
```

### Main Features
- Boid simulation with basic flocking behavior.
- Real-time rendering in the terminal using tcell.
- Configurable parameters for boid behavior (e.g., separation, alignment, cohesion) (in code, i wont be doing args type for configuration yet).
- Predator and prey dynamics.
- Can add new boids/predators using mouse clicks.

### Future
- Performance improvements for larger flocks.
- More system info to display.
- Implement a different rendering option in the terminal for smoother movement and better visuals.

### Contributing
- Feel free to contribute by making a PR~

```
           \\
   \\      (o>
   (o>     //\
___(())____V_/_____ birb
   ||      ||
           ||
```

### Demo
<img width="1358" height="792" alt="boidfetch" src="https://github.com/user-attachments/assets/41f5eaa2-944c-45c9-867d-e72261358d5d" />
