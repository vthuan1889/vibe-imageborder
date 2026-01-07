# Code Review: Phase 2 Template Service

## Scope

- Files: `internal/template/parser.go`, `service.go`, `parser_test.go`, `service_test.go`, `internal/models/types.go`
- LOC: ~260 lines (implementation + tests)
- Review focus: Phase 2 Template Service implementation
- Updated plans: phase-02-template-service.md

## Overall Assessment

Implementation matches plan spec exactly. Code quality solid. Tests pass with 86.5% coverage. No blocking issues.

## Critical Issues

**None**

## High Priority Findings

### H1: Path Traversal Vulnerability (parser.go:20)

```go
func ParseTemplate(path string) (*models.TemplateConfig, error) {
    data, err := os.ReadFile(path)  // ← No path validation
```

**Impact**: User-controlled path → read arbitrary files
**Fix**: Add `filepath.Clean()` + check if path in allowed directory

```go
import "path/filepath"

func ParseTemplate(path string) (*models.TemplateConfig, error) {
    cleanPath := filepath.Clean(path)
    // Optional: validate cleanPath is within allowed directory
    data, err := os.ReadFile(cleanPath)
```

### H2: Race Condition in Cache (service.go:21-31)

```go
func (s *Service) LoadTemplate(path string) (*models.TemplateConfig, error) {
    if cached, ok := s.cache[path]; ok {  // ← Not thread-safe
        return cached, nil
    }
    // ... writes to cache without lock
```

**Impact**: Concurrent access → data races
**Fix**: Add `sync.RWMutex`

```go
type Service struct {
    cache map[string]*models.TemplateConfig
    mu    sync.RWMutex
}

func (s *Service) LoadTemplate(path string) (*models.TemplateConfig, error) {
    s.mu.RLock()
    cached, ok := s.cache[path]
    s.mu.RUnlock()
    if ok {
        return cached, nil
    }

    config, err := ParseTemplate(path)
    if err != nil {
        return nil, err
    }

    s.mu.Lock()
    s.cache[path] = config
    s.mu.Unlock()
    return config, nil
}
```

### H3: Silent strconv.Atoi Error (parser.go:76)

```go
if fs, ok := m["fontsize"].(string); ok {
    size, _ := strconv.Atoi(fs)  // ← Ignores error, defaults to 0
    overlay.FontSize = size
}
```

**Impact**: Invalid fontsize → FontSize=0 → invisible text
**Fix**: Handle error or log

```go
if fs, ok := m["fontsize"].(string); ok {
    size, err := strconv.Atoi(fs)
    if err != nil {
        // Either return error or use default
        size = 12 // default font size
    }
    overlay.FontSize = size
}
```

## Medium Priority Improvements

### M1: ReDoS Potential (parser.go:16)

```go
var fieldRegex = regexp.MustCompile(`\[([^\]]+)\]`)
```

Low risk (simple pattern), but for production consider timeout on `FindAllStringSubmatch`.

### M2: Incomplete Test Coverage (86.5%)

Missing tests:
- Error cases in `parseOverlay` (non-map input)
- Invalid fontsize strings
- Background field edge cases
- ClearCache race conditions

### M3: Memory Leak in Cache

Cache grows unbounded. No eviction policy. For production:
- Add LRU eviction
- Or document cache is per-app-lifecycle

## Low Priority Suggestions

### L1: Code Comments

Add package doc:
```go
// Package template provides template parsing and field extraction
// for JSON-based text overlay templates.
package template
```

### L2: Test Service_test.go Line 40

```go
if config1 != config2 {  // ← Pointer comparison works but fragile
```

Better: `if config1.Background != config2.Background` to test value equality.

## Positive Observations

- Clean separation: parser vs service
- Good error wrapping with `%w`
- YAGNI: No over-engineering
- Tests use `t.TempDir()` correctly
- Regex compiled once (parser.go:16)
- Fields extracted to slice (no map iteration order issues)

## Architecture Violations

**None**. Follows plan architecture exactly.

## YAGNI/KISS/DRY Check

✅ YAGNI: No unused features
✅ KISS: Simple regex, simple cache
✅ DRY: `parseOverlay` extracts common logic

## Security Audit

- ✅ No SQL injection (no DB)
- ❌ **Path traversal** (H1)
- ✅ No XSS (backend only)
- ❌ **Race conditions** (H2)
- ✅ No secrets in code
- ⚠️  User input in `replacePlaceholders` → ensure frontend sanitizes

## Performance Analysis

- Regex compiled once ✅
- Cache prevents re-parsing ✅
- Race detector passes (single-threaded) ⚠️
- No benchmarks ℹ️

**Recommendation**: Add benchmarks if template parsing becomes bottleneck.

## Task Completeness Verification

Plan TODO list (phase-02-template-service.md:375-381):
- [x] Create `internal/template/parser.go` ✅
- [x] Create `internal/template/service.go` ✅
- [x] Create unit tests ✅
- [ ] Test with real template files from reference app ❌
- [ ] Verify field extraction works with complex templates ❌

**Real template tests**: Found fixtures in `tests/fixtures/templates/` but no tests using them.

## Recommended Actions

### Must Fix Before Production

1. **H1**: Add `filepath.Clean()` + path validation in `ParseTemplate`
2. **H2**: Add `sync.RWMutex` to `Service`
3. **H3**: Handle `strconv.Atoi` error

### Should Fix Before Phase 3

4. Add integration tests with real templates from `tests/fixtures/templates/`
5. Increase test coverage to 95%+ (add error cases)

### Nice to Have

6. Add package documentation
7. Add benchmarks

## Metrics

- Type Coverage: N/A (Go)
- Test Coverage: 86.5%
- Linting Issues: 0 (go vet clean)
- Build Status: ✅ Pass
- Race Detector: ✅ Pass (single-threaded test)

## Unresolved Questions

1. Should cache have TTL or size limit?
2. Should invalid fontsize fail parsing or use default?
3. What's the max template file size limit?
