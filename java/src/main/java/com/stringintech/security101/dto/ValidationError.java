package com.stringintech.security101.dto;

import java.time.Instant;
import java.util.List;

public class ValidationError {
    private final String message;
    private final String code;
    private final Instant timestamp;
    private final List<FieldError> fieldErrors;

    public ValidationError(String message, String code, List<FieldError> fieldErrors) {
        this.message = message;
        this.code = code;
        this.timestamp = Instant.now();
        this.fieldErrors = fieldErrors;
    }

    public ValidationError(String message, String code) {
        this(message, code, List.of());
    }

    public String getMessage() {
        return message;
    }

    public String getCode() {
        return code;
    }

    public Instant getTimestamp() {
        return timestamp;
    }

    public List<FieldError> getFieldErrors() {
        return fieldErrors;
    }

    public record FieldError(String field, String message) {
    }
}