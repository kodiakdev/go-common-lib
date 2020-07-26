# Changelog

## Unreleased
### Added
- New err type unique constraint violation
- Common auditing struct
- New constructor and default collection name for DBOperation
- Common repo db operation: FindAtColl, FindOneAtColl, InsertAtColl, UpdateAtColl to enable user specify custom collection
- Suppport for pagination and sort in common repo
- Parser and validator for Basic authentication

## 0.0.1-alpha.3
### Added
- Common repo db operation: Find, FindOne, Insert, Update
- Common response across services such as UnknownError or DBProblem response
- Common error code and explanation across services
- Swagger html error message generator

## 0.0.1-alpha.2
### Fixed
- Module declaration in go.mod

## 0.0.1-alpha.1
### Added
- Check whether error is mongodb's duplication write error