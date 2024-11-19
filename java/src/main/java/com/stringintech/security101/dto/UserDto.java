package com.stringintech.security101.dto;

import com.stringintech.security101.model.User;

import java.time.Instant;

public class UserDto {

    private String fullName;

    private String username;

    private Instant createdAt;

    public UserDto() {
    }

    public UserDto(User user) {
        this.fullName = user.getFullName();
        this.username = user.getUsername();
        this.createdAt = user.getCreatedAt();
    }

    public String getFullName() {
        return fullName;
    }

    public void setFullName(String fullName) {
        this.fullName = fullName;
    }

    public String getUsername() {
        return username;
    }

    public void setUsername(String username) {
        this.username = username;
    }

    public Instant getCreatedAt() {
        return createdAt;
    }

    public void setCreatedAt(Instant createdAt) {
        this.createdAt = createdAt;
    }

}
