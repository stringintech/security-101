package com.stringintech.security101.controller;

import com.stringintech.security101.dto.UserDto;
import com.stringintech.security101.dto.UserLoginRequest;
import com.stringintech.security101.dto.UserLoginResponse;
import com.stringintech.security101.dto.UserRegisterRequest;
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
    public ResponseEntity<UserDto> register(@RequestBody UserRegisterRequest req) {
        User u = this.userService.register(req);

        UserDto dto = new UserDto(u);
        return ResponseEntity.ok(dto);
    }

    @PostMapping("/auth/login")
    public ResponseEntity<UserLoginResponse> login(@RequestBody UserLoginRequest req) {
        User u = this.userService.authenticate(req);
        Instant creationTime = timeService.now();
        String token = jwtService.generateToken(u, creationTime);

        UserLoginResponse resp = new UserLoginResponse();
        resp.setUser(u);
        resp.setToken(token);
        resp.setTokenCreationTime(creationTime);
        return ResponseEntity.ok(resp);
    }

    @PostMapping("/users/me")
    public ResponseEntity<UserDto> getAuthenticatedUser() {
        Authentication auth = SecurityContextHolder.getContext().getAuthentication();
        User u = (User) auth.getPrincipal();
        UserDto dto = new UserDto(u);
        return ResponseEntity.ok(dto);
    }
}