package repository

import "github.com/fikryfahrezy/let-it-go/pkg/app_error"

// Repository errors
var (
	// User not found errors
	ErrUserNotFound = app_error.New("USER-USER_NOT_FOUND", "user not found")

	// Database operation errors
	ErrFailedToCreateUser     = app_error.New("USER-FAILED_TO_CREATE_USER", "failed to create user")
	ErrFailedToGetUser        = app_error.New("USER-FAILED_TO_GET_USER", "failed to get user")
	ErrFailedToGetUserByEmail = app_error.New("USER-FAILED_TO_GET_USER_BY_EMAIL", "failed to get user by email")
	ErrFailedToUpdateUser     = app_error.New("USER-FAILED_TO_UPDATE_USER", "failed to update user")
	ErrFailedToDeleteUser     = app_error.New("USER-FAILED_TO_DELETE_USER", "failed to delete user")
	ErrFailedToListUsers      = app_error.New("USER-FAILED_TO_LIST_USERS", "failed to list users")
	ErrFailedToCountUsers     = app_error.New("USER-FAILED_TO_COUNT_USERS", "failed to count users")

	// Row scanning errors
	ErrFailedToScanUserRow = app_error.New("USER-FAILED_TO_SCAN_USER_ROW", "failed to scan user row")

	// Database result errors
	ErrFailedToGetLastInsertID = app_error.New("USER-FAILED_TO_GET_LAST_INSERT_ID", "failed to get last insert id")
	ErrFailedToGetRowsAffected = app_error.New("USER-FAILED_TO_GET_ROWS_AFFECTED", "failed to get rows affected")
	ErrFailedToIterateRows     = app_error.New("USER-FAILED_TO_ITERATE_ROWS", "error iterating rows")
)