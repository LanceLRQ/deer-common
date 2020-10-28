package structs

// 题目Input/Output样例信息
type ProblemIOSample struct {
    Input  string `json:"input"`  // Input sample
    Output string `json:"output"` // Output sample
}

// 题目正文信息
type ProblemContent struct {
    Author      string            `json:"author"`      // Problem author
    Source      string            `json:"source"`      // Problem source
    Description string            `json:"description"` // Description
    Input       string            `json:"input"`       // Input requirements
    Output      string            `json:"output"`      // Output requirements
    Sample      []ProblemIOSample `json:"sample"`      // Sample cases
    Tips        string            `json:"tips"`        // Solution tips
}
