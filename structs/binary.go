package structs

type ShellResult struct {
    Success     bool
    Stdout      string
    Stderr      string
    ExitCode    int
}
