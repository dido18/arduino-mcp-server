package main

import (
	"context"
	"errors"
	"fmt"
	"os/exec"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	s := server.NewMCPServer(
		"Arduino MCP",
		"0.0.1",
	)

	listboardsTool := mcp.NewTool("list_boards",
		mcp.WithDescription("Detects and displays a list of boards connected to the current computer."),
	)

	uploadTool := mcp.NewTool("upload",
		mcp.WithDescription("Upload Sketch to the board"),
		mcp.WithString("fqbn",
			mcp.Required(),
			mcp.Description("Fully Qualified Board Name, e.g.: arduino:avr:uno"),
		),
		mcp.WithString("port",
			mcp.Required(),
			mcp.Description("Upload port address, e.g.: COM3 or /dev/ttyACM2"),
		),
		mcp.WithString("sketch",
			mcp.Required(),
			mcp.Description("Path to the sketch file, e.g.: /path/to/Blink.ino"),
		),
	)

	compileTool := mcp.NewTool("compile",
		mcp.WithDescription("Compile (and upload) Arduino Sketch"),
		mcp.WithString("fqbn",
			mcp.Required(),
			mcp.Description("Fully Qualified Board Name, e.g.: arduino:avr:uno"),
		),
		mcp.WithString("sketch",
			mcp.Required(),
			mcp.Description("Path to the sketch file, e.g.: /path/to/Blink.ino"),
		),
		mcp.WithBoolean("upload",
			mcp.Description("Upload the sketch after compiling"),
			mcp.DefaultBool(false),
		),
		mcp.WithString("port",
			mcp.Description("Upload port address, e.g.: COM3 or /dev/ttyACM2"),
		),
	)

	s.AddTool(listboardsTool, listBoardsHandler)
	s.AddTool(uploadTool, uploadHandler)
	s.AddTool(compileTool, compileHandler)

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

func listBoardsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	cmd := exec.Command("arduino-cli", "board", "list", "--json")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute arduino-cli: %w", err)
	}

	return mcp.NewToolResultText(string(output)), nil
}

func uploadHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	fqbn, ok := request.Params.Arguments["fqbn"].(string)
	if !ok {
		return nil, errors.New("fqbn must be a string")
	}
	port, ok := request.Params.Arguments["port"].(string)
	if !ok {
		return nil, errors.New("port must be a string")
	}
	sketch, ok := request.Params.Arguments["sketch"].(string)
	if !ok {
		return nil, errors.New("sketch must be a string")
	}

	args := []string{"upload", "--json", "--fqbn", fqbn, "--port", port, sketch}

	cmd := exec.Command("arduino-cli", args...)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute arduino-cli: %w", err)
	}

	return mcp.NewToolResultText(string(output)), nil
}

func compileHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	fqbn, ok := request.Params.Arguments["fqbn"].(string)
	if !ok {
		return nil, errors.New("fqbn must be a string")
	}
	sketch, ok := request.Params.Arguments["sketch"].(string)
	if !ok {
		return nil, errors.New("sketch must be a string")
	}

	args := []string{"compile", "--json", "--fqbn", fqbn, sketch}

	upload, ok := request.Params.Arguments["upload"].(bool)
	if ok && upload {
		args = append(args, "--upload")
	}
	if upload {
		port, ok := request.Params.Arguments["port"].(string)
		if !ok {
			return nil, errors.New("port must be a string")
		}
		args = append(args, "--port", port)
	}

	cmd := exec.Command("arduino-cli", args...)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute arduino-cli: %w", err)
	}

	return mcp.NewToolResultText(string(output)), nil
}
