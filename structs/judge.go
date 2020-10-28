package structs

// 特殊评测设置
type SpecialJudgeOptions struct {
    Mode               int    `json:"mode"`                 // Mode；0-Disabled；1-Normal；2-Interactor
    Checker            string `json:"checker"`              // Checker file path (Use code file is better then compiled binary!)
    RedirectProgramOut bool   `json:"redirect_program_out"` // Redirect target program's STDOUT to checker's STDIN (checker mode). if not, redirect testcase-in file to checker's STDIN
    TimeLimit          int    `json:"time_limit"`           // Time limit (ms)
    MemoryLimit        int    `json:"memory_limit"`         // Memory limit (kb)
}

// 测试数据
type TestCase struct {
    Handle      string `json:"handle"`              // Identifier
    Name        string `json:"name"`                // Testcase name
    TestCaseIn  string `json:"test_case_in"`        // Testcase input file path
    TestCaseOut string `json:"test_case_out"`       // Testcase output file path
}

// 测试数据运行结果
type TestCaseResult struct {
    Handle       string `json:"handle"`             // Identifier

    TestCaseIn   string `json:"-"`                  // Testcase input file path (internal)
    TestCaseOut  string `json:"-"`                  // Testcase output file path (internal)
    ProgramOut   string `json:"program_out"`        // Program-stdout file path
    ProgramError string `json:"program_error"`      // Program-stderr file path

    JudgerOut    string `json:"judger_out"`         // Special judger checker's stdout
    JudgerError  string `json:"judger_error"`       // Special judger checker's stderr
    JudgerReport string `json:"judger_report"`      // Special judger checker's report file

    JudgeResult int    `json:"judge_result"`        // Judge result flag number
    TextDiffLog string `json:"text_diff_log"`       // Text Checkup Log
    TimeUsed    int    `json:"time_used"`           // Maximum time used
    MemoryUsed  int    `json:"memory_used"`         // Maximum memory used
    ReSignum    int    `json:"re_signal_num"`       // Runtime error signal number
    SameLines   int    `json:"same_lines"`          // Same lines when WA
    TotalLines  int    `json:"total_lines"`         // Total lines when WA
    ReInfo      string `json:"re_info"`             // ReInfo when Runtime Error or special judge Runtime Error
    SeInfo      string `json:"se_info"`             // SeInfo when System Error
    CeInfo      string `json:"ce_info"`             // CeInfo when Compile Error

    SPJExitCode   int `json:"spj_exit_code"`        // Special judge exit code
    SPJTimeUsed   int `json:"spj_time_used"`        // Special judge maximum time used
    SPJMemoryUsed int `json:"spj_memory_used"`      // Special judge maximum memory used
    SPJReSignum   int `json:"spj_re_signal_num"`    // Special judge runtime error signal number
}

// 评测结果信息
type JudgeResult struct {
    SessionId   string           `json:"session_id"`    // Judge Session Id
    JudgeResult int              `json:"judge_result"`  // Judge result flag number
    TimeUsed    int              `json:"time_used"`     // Maximum time used
    MemoryUsed  int              `json:"memory_used"`   // Maximum memory used
    TestCases   []TestCaseResult `json:"test_cases"`    // Testcase Results
    ReInfo      string           `json:"re_info"`       // ReInfo when Runtime Error or special judge Runtime Error
    SeInfo      string           `json:"se_info"`       // SeInfo when System Error
    CeInfo      string           `json:"ce_info"`       // CeInfo when Compile Error
}

// 评测资源限制信息
type JudgeResourceLimit struct {
    TimeLimit     int `json:"time_limit"`               // Time limit (ms)
    MemoryLimit   int `json:"memory_limit"`             // Memory limit (KB)
    RealTimeLimit int `json:"real_time_limit"`          // Real Time Limit (ms) (optional)
    FileSizeLimit int `json:"file_size_limit"`          // File Size Limit (bytes) (optional)
}

// 评测配置信息
type JudgeConfiguration struct {
    TestCases     []TestCase                    `json:"test_cases"`      // Test cases
    TimeLimit     int                           `json:"time_limit"`      // Time limit (ms)
    MemoryLimit   int                           `json:"memory_limit"`    // Memory limit (KB)
    RealTimeLimit int                           `json:"real_time_limit"` // Real Time Limit (ms) (optional)
    FileSizeLimit int                           `json:"file_size_limit"` // File Size Limit (bytes) (optional)
    Uid           int                           `json:"uid"`             // User id (optional)
    StrictMode    bool                          `json:"strict_mode"`     // Strict Mode (if close, PE will be ignore)
    SpecialJudge  SpecialJudgeOptions           `json:"special_judge"`   // Special Judge Options
    Limitation    map[string]JudgeResourceLimit `json:"limitation"`      // Limitation
    Problem       ProblemContent                `json:"problem"`         // Problem Info
}
