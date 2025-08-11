
# Use Case Scenarios

List of use cases and all usage related options needed to understand intent and full feature set
of Packager.

## Build Image

Build image based on Dockerfile in Context. It is used before `build-package` command.

### Build Single Image

Build single image specified by name.

**Command**

```bash
packager build-image
  --context ./example \
  --name debian
```

### Build All Images

Build all images in Context.

**Command**

```bash
packager build-image
  --context ./example \
  --all
```

## Build Package

### Dependencies

Each use case is described by a simple mermaid diagram which describes dependency graph of noted
Packages.

Green color indicates which Packages will be built.

Blue colored arrows indicate which Packages are built because of `--build-deps-on` or
`--build-deps-on-recursive` option.

Arrows indicate dependency (build deps) between Packages. When Package A depends on Package B the
following is written.

```mermaid
graph TD
    A
    A --> B
```

### Build Package - without Dependencies

Build single Package (G) without any dependencies.

It expects all Package dependencies are already built and installed into build sysroot directory.

```mermaid
graph TD
    A ---> B
    A --> D
    B --> E
    B --> F
    C --> G
    D --> G
    H --> L
    D --> I
    E --> J
    F --> K
    G --> L
    I --> M
    L --> N
    L --> O
    P --> Q

    style G color:green;
```

**Command**

```bash
packager build-package
  --context ./example \
  --image-name debian \
  --name G \
  --output-dir ./git-lfs-repo
```

### Build Package - with Dependencies

Build all dependencies (L, N, O) of the Package (G) before the Package (G) is built.

```mermaid
graph TD
    A ---> B
    A --> D
    B --> E
    B --> F
    C --> G
    D --> G
    H --> L
    D --> I
    E --> J
    F --> K
    G --> L
    I --> M
    L --> N
    L --> O
    P --> Q

    style G color:green;
    style L color:green;
    style N color:green;
    style O color:green;
    linkStyle 10 stroke:green,stroke-width:2px;
    linkStyle 12 stroke:green,stroke-width:2px;
    linkStyle 13 stroke:green,stroke-width:2px;
```

**Command**

```bash
packager build-package
  --context ./example \
  --image-name debian \
  --name G \
  --build-deps \
  --output-dir ./git-lfs-repo
```

**WARNING**: Building only some Packages can break dependency relationships. In this example,
Packages G, L, N, and O are built, but other Packages that depend on them (C, D, H) are not rebuilt.
If the build environment changes (e.g., compiler version, system libraries), the previously built
Packages C, D, and H may become incompatible with the newly built dependencies. While this scenario
is uncommon, it can cause runtime failures. Currently, Packager does not detect these issues
automatically. 

### Build Package - with Depends on Packages

Build Packages (C, D) which depends on Package (G) with its dependencies (I, M) without Package (G)
and its dependencies (L, N, O).

```mermaid
graph TD
    A ---> B
    A --> D
    B --> E
    B --> F
    C --> G
    D --> G
    H --> L
    D --> I
    E --> J
    F --> K
    G --> L
    I --> M
    L --> N
    L --> O
    P --> Q

    style C color:green;
    style D color:green;
    style I color:green;
    style M color:green;
    linkStyle 4 stroke:royalblue;
    linkStyle 5 stroke:royalblue;
    linkStyle 7 stroke:green,stroke-width:2px;
    linkStyle 11 stroke:green,stroke-width:2px;
```

**Command**

```bash
packager build-package
  --context ./example \
  --image-name debian \
  --name G \
  --build-deps-on \
  --output-dir ./git-lfs-repo
```

> **NOTE**: The `--build-deps` option can be added to command to build also the G Package and its
dependencies (L, N, O).

### Build Package - with Depends on Packages Recursive

Build Packages (C, D, A) which depends on Package (G) recursively with its dependencies
(B, D, E, F, I, J, K, M) without Package (G) and its dependencies (L, N, O).

```mermaid
graph TD
    A ---> B
    A --> D
    B --> E
    B --> F
    C --> G
    D --> G
    H --> L
    D --> I
    E --> J
    F --> K
    G --> L
    I --> M
    L --> N
    L --> O
    P --> Q

    style C color:green;
    style D color:green;
    style I color:green;
    style M color:green;
    style A color:green;
    style B color:green;
    style E color:green;
    style F color:green;
    style J color:green;
    style K color:green;
    linkStyle 0 stroke:green,stroke-width:2px;
    linkStyle 1 stroke:green,stroke-width:2px;
    linkStyle 2 stroke:green,stroke-width:2px;
    linkStyle 3 stroke:green,stroke-width:2px;
    linkStyle 4 stroke:royalblue;
    linkStyle 5 stroke:royalblue;
    linkStyle 7 stroke:green,stroke-width:2px;
    linkStyle 8 stroke:green,stroke-width:2px;
    linkStyle 9 stroke:green,stroke-width:2px;
    linkStyle 11 stroke:green,stroke-width:2px;
```

**Command**

```bash
packager build-package
  --context ./example \
  --image-name debian \
  --name G \
  --build-deps-on-recursive \
  --output-dir ./git-lfs-repo
```

> **NOTE**: The `--build-deps` option can be added to command to build also the G Package and its
dependencies (L, N, O).

### Build Package - all Packages

Build all Packages in Context.

```mermaid
graph TD
    A ---> B
    A --> D
    B --> E
    B --> F
    C --> G
    D --> G
    H --> L
    D --> I
    E --> J
    F --> K
    G --> L
    I --> M
    L --> N
    L --> O
    P --> Q

    style A color:green;
    style B color:green;
    style C color:green;
    style D color:green;
    style E color:green;
    style F color:green;
    style G color:green;
    style H color:green;
    style I color:green;
    style J color:green;
    style K color:green;
    style L color:green;
    style M color:green;
    style N color:green;
    style O color:green;
    style P color:green;
    style Q color:green;
    linkStyle 0 stroke:green,stroke-width:2px;
    linkStyle 1 stroke:green,stroke-width:2px;
    linkStyle 2 stroke:green,stroke-width:2px;
    linkStyle 3 stroke:green,stroke-width:2px;
    linkStyle 4 stroke:green,stroke-width:2px;
    linkStyle 5 stroke:green,stroke-width:2px;
    linkStyle 6 stroke:green,stroke-width:2px;
    linkStyle 7 stroke:green,stroke-width:2px;
    linkStyle 8 stroke:green,stroke-width:2px;
    linkStyle 9 stroke:green,stroke-width:2px;
    linkStyle 10 stroke:green,stroke-width:2px;
    linkStyle 11 stroke:green,stroke-width:2px;
    linkStyle 12 stroke:green,stroke-width:2px;
    linkStyle 13 stroke:green,stroke-width:2px;
    linkStyle 14 stroke:green,stroke-width:2px;
```

**Command**

```bash
packager build-package
  --context ./example \
  --image-name debian \
  --all \
  --output-dir ./git-lfs-repo
```

## Build App

For apps which manages their own dependencies (e. g. using CMake), the Packager supports different
workflow - Apps. The Configs for these Apps are in `app` directory in Context and are the same as
for regular Packages with one exception - DependsOn array must be empty. The Apps do not support
any dependencies managed by Packager, if the user needs it, the Packages should be used. Building
Apps also do not support any dependency related flags (`--build-deps`, `--build-deps-on` ...).

### Build App - all Apps

Build all Apps in Context.

**Command**

```bash
packager build-app
  --context ./example \
  --image-name debian \
  --all \
  --output-dir ./git-lfs-repo
```

### Build App - single App

Build single App in Context.

**Command**

```bash
packager build-app
  --context ./example \
  --image-name debian \
  --name app-name \
  --output-dir ./git-lfs-repo
```

## Create Sysroot

When all Packages are build and stored as part of `--output-dir` directory the sysroot can be
created.

The packager takes all archives for a given Image name and Architecture and unzip them into the
specified directory.

**Command**

Creates sysroot from Packages in Package Repository for `debian` image in `new_sysroot/` directory.

```bash
packager create-sysroot
  --context ./example \
  --image-name debian \
  --git-lfs ./git-lfs-repo \
  --sysroot-dir new_sysroot
```
