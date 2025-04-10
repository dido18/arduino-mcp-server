# arduino-mcp-server
An arduino mcp server written in golang.

## Installation

### Vscode
- `git clone git@github.com:dido18/arduino-mcp-server.git`
- `cd arduino-mcp-server && go build .`
- move the binary into a location (e.g. `mv arduino-mcp-server ~/arduino-mcp-server` )
- use the location in the `command` field of the `settings.json` (or can add it to a file called `.vscode/mcp.json` in your workspace) 

```json
{
    // settings.json
    "mcp": {
        "servers": {
            "arduino-mcp-server": {
                "type": "stdio",
                "command": "~/arduino-mcp-server",
                "args": []
            }   
        }
    },
}

```

## Tools:

- `list_boards`: list the connected board to the pc
- `compile`: compile (and optionally upload) a sketch
    - `fqbn` the fqbn of the board
    - `sketch` the path of the sketch to compile 
    - `upload` if True perform the upload
    - `port` the port wher to upload
- `upload`: upload 
    - `fqbn` the fqbn of the board
    - `sketch` the path of the sketch to compile 
    - `port` the port where to upload
