package com.stringintech.security101.dto;

import com.stringintech.security101.model.User;

import java.time.Instant;

public class UserLoginResponse {

    private User user;

    private String token;

    private Instant tokenCreationTime;

    public User getUser() {
        return user;
    }

    public void setUser(User user) {
        this.user = user;
    }

    public String getToken() {
        return token;
    }

    public void setToken(String token) {
        this.token = token;
    }

    public Instant getTokenCreationTime() {
        return tokenCreationTime;
    }

    public void setTokenCreationTime(Instant tokenCreationTime) {
        this.tokenCreationTime = tokenCreationTime;
    }
}
