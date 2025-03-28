# Changelog

## [Unreleased]

### Changes
- Updated minimum Go version to 1.20
- Renamed `HttpError` to `HTTPError` to follow Go naming conventions (with backward compatibility aliases)
- Renamed `ParseDataUri` to `ParseDataURI` for consistent naming (with backward compatibility function)
- Added backward compatibility for renamed functions while marking old names as deprecated
- Added proper error wrapping with `%w` verb in `fmt.Errorf`
- Added context support in HTTP operations
- Added named error variables for common error cases
- Improved error handling with `errors.Is` and `errors.As` instead of type assertions
- Enhanced documentation with detailed function descriptions and examples
- Added nil checks and defensive programming patterns
- Refactored the PHP query string parser for better reliability and maintainability
- Made test cases more descriptive with subtests and proper error messages
- Updated README with feature summary and usage examples

### Bug fixes
- Fixed edge cases in resumable downloads
- Properly handle EOF conditions in Read operations
- Added safety checks to prevent nil pointer dereferences
- Fixed handling of empty host in ParseIPPort

### Performance improvements
- Optimized string handling in ParseDataURI using TrimLeft instead of byte-by-byte removal
- Improved memory usage in resumeget.go with proper resource cleanup