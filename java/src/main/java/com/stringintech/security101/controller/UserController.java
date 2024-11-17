package com.stringintech.security101.controller;

import com.stringintech.security101.dto.UserLoginDto;
import com.stringintech.security101.dto.UserRegisterDto;
import com.stringintech.security101.model.User;
import com.stringintech.security101.service.JwtService;
import com.stringintech.security101.service.TimeService;
import com.stringintech.security101.service.UserService;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RestController;

import java.time.Instant;


@RestController
public class UserController {

    private final UserService userService;
    private final JwtService jwtService;
    private final TimeService timeService;

    public UserController(UserService userService, JwtService jwtService, TimeService timeService) {
        this.userService = userService;
        this.jwtService = jwtService;
        this.timeService = timeService;
    }

    @PostMapping("/auth/register")
    public ResponseEntity<User> register(@RequestBody UserRegisterDto user) { //TODO validation return type
        User u = this.userService.register(user);
        return ResponseEntity.ok(u); //TODO handle exceptions
    }

    @PostMapping("/auth/login")
    public ResponseEntity<UserLoginResponse> login(@RequestBody UserLoginDto user) {
        User u = this.userService.authenticate(user);
        Instant creationTime = timeService.now();
        String token = jwtService.generateToken(u, creationTime);
        UserLoginResponse r = new UserLoginResponse();
        r.setUser(u);
        r.setToken(token);
        r.setTokenCreationTime(creationTime);
        return ResponseEntity.ok(r);
    }

    @PostMapping("/users/me")
    public ResponseEntity<User> getAuthenticatedUser() {
        Authentication auth = SecurityContextHolder.getContext().getAuthentication();
        User user = (User) auth.getPrincipal();
        return ResponseEntity.ok(user);
    }
}