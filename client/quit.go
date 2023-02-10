package client

type quitArgs struct{}

func init() {
	parser.AddCommand("q", "quit the program", "", &quitArgs{})
}

func (q *quitArgs) Execute(args []string) error {
	client.shouldExit = true
	return nil
}
