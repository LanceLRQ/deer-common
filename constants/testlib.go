package constants

var TestlibBinaryPrefixs = map[string]string {
    "generator": "g_",
    "validator": "",
    "checker": "",
    "interactor": "",
}

var TestlibExitCodeMapping = map[int]int {
    0: JudgeFlagAC,
    1: JudgeFlagWA,
    2: JudgeFlagPE,
    3: JudgeFlagSpecialJudgeError,
    4: JudgeFlagSpecialJudgeError,
    5: JudgeFlagSpecialJudgeError,
    8: JudgeFlagSpecialJudgeError,
    16: JudgeFlagWA,
}
