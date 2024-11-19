package com.stringintech.security101.repository;

import com.stringintech.security101.exception.DuplicateUsernameException;
import com.stringintech.security101.model.User;
import org.springframework.stereotype.Service;

import java.util.concurrent.ConcurrentHashMap;

@Service
public class UserRepository {

    private final ConcurrentHashMap<String, User> users = new ConcurrentHashMap<>();

    public UserRepository() {
    }

    public User createUser(User user) {
        User existing = users.putIfAbsent(user.getUsername(), user);
        if (existing != null) {
            throw new DuplicateUsernameException(user.getUsername());
        }
        return user;
    }

    public User getUserByUsername(String username) {
        return users.get(username);
    }
}
