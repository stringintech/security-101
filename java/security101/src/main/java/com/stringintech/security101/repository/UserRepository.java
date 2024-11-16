package com.stringintech.security101.repository;

import com.stringintech.security101.model.User;
import org.springframework.stereotype.Service;

import java.util.concurrent.ConcurrentHashMap;

@Service
public class UserRepository {

    private final ConcurrentHashMap<String, User> users = new ConcurrentHashMap<>();

    public UserRepository() {
    }

    public User createUser(User user) {
        if (users.containsKey(user.getUsername())) {
            throw new IllegalArgumentException("Duplicate username"); //TODO or invalid request? enhance exception handling
        }
        users.put(user.getUsername(), user);
        return user;
    }

    public User getUserByUsername(String username) {
        return users.get(username);
    }
}
