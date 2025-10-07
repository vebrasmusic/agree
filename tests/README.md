# Parser Tests

Comprehensive test suite for the `agree` parser that extracts type information from Python code.

## Test Coverage

### 1. Pydantic Schemas (`TestPydanticSchemas`)
- **Simple types**: Basic type annotations (int, str, float, bool)
- **Optional types**: `Optional[T]` handling
- **Union types**: `Union[T1, T2, ...]` with multiple types
- **Pipe unions**: Modern `T1 | T2` syntax
- **Nested unions**: Complex combinations like `Optional[Union[int, str]]`

### 2. SQLAlchemy New Style (`TestSQLAlchemyNewStyle`)
- **Mapped types**: `Mapped[int]`, `Mapped[str]`, etc.
- **Optional Mapped**: `Mapped[Optional[T]]` handling
- **DateTime types**: datetime, date, time support

### 3. SQLAlchemy Old Style (`TestSQLAlchemyOldStyle`)
- **Column() calls**: `Column(Integer)`, `Column(String)`, etc.
- **Nullable columns**: `nullable=True` parameter detection
- **Type mappings**: Integer→int, String→str, DateTime→datetime, etc.
- **Special handling**: Skips `__tablename__` and other dunder attributes

### 4. Mixed Styles (`TestMixedStyles`)
- **Multiple classes, same target**: Pydantic + SQLAlchemy old + new targeting same entity
- **Multiple targets**: Different classes targeting different entities
- **Consistency**: Verifies all styles produce identical output for equivalent schemas

### 5. Edge Cases (`TestEdgeCases`)
- **Non-decorated classes**: Classes without `@agree` are ignored
- **Empty classes**: Classes with no fields
- **Large unions**: 6+ types in a single union
- **Deep nesting**: Multiple levels of Optional/Union nesting

### 6. Decorator Parameters (`TestAgreeDecorator`)
- **Target extraction**: `@agree(target="...")`
- **Multiple parameters**: `@agree(target="...", fidelity=2)`

## Running Tests

Run all tests:
```bash
uv run pytest tests/
```

Run with verbose output:
```bash
uv run pytest tests/ -v
```

Run specific test class:
```bash
uv run pytest tests/test_parser.py::TestPydanticSchemas -v
```

Run specific test:
```bash
uv run pytest tests/test_parser.py::TestPydanticSchemas::test_simple_pydantic_schema -v
```

## Test Statistics

- **Total tests**: 20
- **Test classes**: 6
- **Coverage**: All three supported styles (Pydantic, SQLAlchemy new, SQLAlchemy old)

## Adding New Tests

When adding new tests:

1. Choose the appropriate test class based on what you're testing
2. Follow the naming convention: `test_<feature_description>`
3. Include a docstring explaining what the test validates
4. Use triple-quoted strings for test code
5. Assert both structure and values in the result

Example:
```python
def test_my_new_feature(self):
    """Test description here"""
    code = '''
from pydantic import BaseModel

@agree(target="MyTarget")
class MySchema(BaseModel):
    field: int
'''
    result = parse_code(code)
    assert result["MyTarget"]["MySchema"]["fields"]["field"] == ["int"]
```
