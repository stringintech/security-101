package com.stringintech.security101.controller;

import com.stringintech.security101.model.User;
import com.stringintech.security101.service.UserService;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class UserController {

    private final UserService userService;

    public UserController(UserService userService) {
        this.userService = userService;
    }

    @PostMapping("/auth/create")
    public User createUser(@RequestBody User user) { //TODO validation return type
        return this.userService.createUser(user); //TODO handle exceptions
    }
}