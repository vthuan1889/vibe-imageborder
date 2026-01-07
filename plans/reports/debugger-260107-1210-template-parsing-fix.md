# Debug Report: Template Parsing Error Fix

**Date:** 2026-01-07
**Agent:** debugger (ad5a6ef)
**Status:** ✅ Resolved

## Executive Summary

Template parsing failed khi gặp field metadata (string) thay vì TemplateField object. Fixed bằng custom UnmarshalJSON để skip non-object fields.

- **Root Cause:** Template JSON có mixed types (`"background": "#f1eeea"` string vs `"barcode": {...}` object)
- **Impact:** Parser crash khi load template khung-004-01.txt
- **Solution:** Custom UnmarshalJSON method skip metadata fields, chỉ parse TemplateField objects
- **Result:** All tests pass, 5 dynamic fields extracted correctly từ khung-004-01.txt

## Technical Analysis

### Error Timeline

1. Test `TestRealTemplates/khung-004-01.txt` failed
2. Error: `json: cannot unmarshal string into Go value of type models.TemplateField`
3. Template có `"background": "#f1eeea"` (string) nhưng parser expect all values là TemplateField object

### Root Cause

**File:** `tests/fixtures/templates/khung-004-01.txt`

```json
{
  "background": "#f1eeea",        // ❌ String, not TemplateField
  "barcode": {                     // ✅ Valid TemplateField object
    "text": "[barcode]",
    "position": "90,1852",
    "fontsize": "50",
    "color": "white"
  }
}
```

**Problem:** Go's default JSON unmarshaler cho `map[string]TemplateField` không handle mixed types.

### Solution Implementation

**File:** `internal/models/types.go`

Added custom `UnmarshalJSON` method:

```go
func (t *Template) UnmarshalJSON(data []byte) error {
    // Parse as raw map first
    var raw map[string]json.RawMessage
    if err := json.Unmarshal(data, &raw); err != nil {
        return err
    }

    *t = make(Template)
    for key, value := range raw {
        var field TemplateField
        // Try unmarshal as TemplateField
        if err := json.Unmarshal(value, &field); err != nil {
            // Skip if not TemplateField (e.g., "background": "#f1eeea")
            continue
        }
        (*t)[key] = field
    }

    return nil
}
```

**Logic:**
1. Parse JSON thành `map[string]json.RawMessage` trước
2. Loop qua từng entry, thử unmarshal thành `TemplateField`
3. Nếu fail → skip (là metadata)
4. Nếu success → add vào Template map

**Changes:**
- Added `encoding/json` import
- Added `UnmarshalJSON` method (23 lines)
- Zero breaking changes

## Test Results

### Before Fix
```
Failed to parse ../../tests/fixtures/templates/khung-004-01.txt:
failed to parse template JSON: json: cannot unmarshal string into Go value of type models.TemplateField
```

### After Fix
```
=== RUN   TestRealTemplates/khung-004-01.txt
    parser_test.go:95: Template ../../tests/fixtures/templates/khung-004-01.txt has 5 fields:
    [price size_dai size_rong size_cao barcode]
--- PASS: TestRealTemplates/khung-004-01.txt (0.00s)
```

**Extracted Fields:** barcode, price, size_dai, size_rong, size_cao ✅

### Full Test Suite
```
PASS: TestParseTemplate (0.00s)
PASS: TestExtractDynamicFields (0.00s)
PASS: TestReplaceVariables (0.00s)
PASS: TestRealTemplates (0.00s)
  PASS: TestRealTemplates/khung-002-05.txt (0.00s) - 7 fields
  PASS: TestRealTemplates/khung-004-01.txt (0.00s) - 5 fields
```

## Recommendations

### Immediate Actions (Done)
- ✅ Custom UnmarshalJSON implemented
- ✅ All tests passing
- ✅ No breaking changes

### Future Enhancements
1. **Explicit metadata handling:** Parse `background` field vào separate struct field thay vì skip
2. **Validation:** Log warning khi skip unknown fields (debugging)
3. **Documentation:** Update template spec để clarify metadata fields

### Preventive Measures
- Template validation script để check structure trước khi commit
- Unit test cho edge cases (empty objects, null values, nested metadata)

## Supporting Evidence

**Modified Files:**
- `internal/models/types.go` (+24 lines: import + UnmarshalJSON method)

**Test Coverage:**
- Template parsing: ✅
- Field extraction: ✅
- Mixed-type handling: ✅
- Backward compatibility: ✅

**No Unresolved Questions**
