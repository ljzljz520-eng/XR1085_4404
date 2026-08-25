# BUG_REPRO

The following failures were observed while validating the initial project state.
Each section records what failed, how to reproduce it, and the complete command output.
They are preserved intentionally; only failing build gates are omitted from the generated Dockerfile.

## Failure 1: Go test (.)

- Observed problem: `Go test (.)` failed in the initial project state.
- Working directory: `.`
- Command: `cd /app && GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go test -count=1 ./...`
- Exit status: `1`

```text
--- FAIL: TestCandleAnimationReturnsToSingleIdleState (0.00s)
    animation_regression_test.go:26: expected merged click layer, got 2
FAIL
FAIL	memorialcandle	0.003s
ok  	memorialcandle/animation	0.002s
?   	memorialcandle/cmd/candle-server	[no test files]
ok  	memorialcandle/cli	0.014s
ok  	memorialcandle/domain	0.001s
ok  	memorialcandle/render	0.002s
ok  	memorialcandle/service	0.027s
ok  	memorialcandle/store	0.016s
ok  	memorialcandle/validate	0.001s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/candle-server): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/candle-server): exit `0`
