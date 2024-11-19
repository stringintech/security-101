package com.stringintech.security101.controller;

import com.stringintech.security101.dto.ValidationError;
import com.stringintech.security101.exception.DuplicateUsernameException;
import com.stringintech.security101.exception.UserRegistrationException;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.MethodArgumentNotValidException;
import org.springframework.web.bind.annotation.ControllerAdvice;
import org.springframework.web.bind.annotation.ExceptionHandler;

import java.util.List;
import java.util.stream.Collectors;

@ControllerAdvice
public class GlobalExceptionHandler {

    @ExceptionHandler(MethodArgumentNotValidException.class)
    public ResponseEntity<ValidationError> handleValidationExceptions(MethodArgumentNotValidException ex) {
        List<ValidationError.FieldError> fieldErrors = ex.getBindingResult()
                .getFieldErrors()
                .stream()
                .map(error -> new ValidationError.FieldError(
                        error.getField(),
                        error.getDefaultMessage()))
                .collect(Collectors.toList());

        ValidationError error = new ValidationError(
                "Validation failed",
                "VALIDATION_ERROR",
                fieldErrors
        );

        return new ResponseEntity<>(error, HttpStatus.BAD_REQUEST);
    }

    @ExceptionHandler(DuplicateUsernameException.class)
    public ResponseEntity<ValidationError> handleDuplicateUsername(DuplicateUsernameException ex) {
        ValidationError error = new ValidationError(
                ex.getMessage(),
                "USER_DUPLICATE_USERNAME"
        );
        return new ResponseEntity<>(error, HttpStatus.BAD_REQUEST);
    }

    @ExceptionHandler(UserRegistrationException.class)
    public ResponseEntity<ValidationError> handleUserRegistration(UserRegistrationException ex) {
        ValidationError error = new ValidationError(
                ex.getMessage(),
                "USER_REGISTRATION_ERROR"
        );
        return new ResponseEntity<>(error, HttpStatus.INTERNAL_SERVER_ERROR); // internal error makes more sense here
    }
}
