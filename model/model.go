package model

type CommandRequest struct {
	Command string `json:"command"`
}

type CommandResponse struct {
	Output string `json:"output,omitempty"`
	Error  string `json:"error,omitempty"`
}

type CommandStatus struct {
	Output string
	Error  string
	Done   bool
}
