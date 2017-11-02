package sesha3

type TokenPayload struct {
	Username string `json:"username"`
	Passwd   string `json:"passwd"`
}

type ExecScriptInstanceResponse struct {
	StackId string `json:"stack_id"`
	Ip      string `json:"ip"`
	CmdOut  []byte `json:"out"`
}

type ExecScriptResponse struct {
	Outputs []ExecScriptInstanceResponse `json:"outputs"`
}
