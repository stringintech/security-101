package com.stringintech.security101.service;

import com.stringintech.security101.dto.UserLoginDto;
import com.stringintech.security101.dto.UserRegisterDto;
import com.stringintech.security101.model.User;
import com.stringintech.security101.repository.UserRepository;
import org.springframework.security.authentication.AuthenticationManager;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;

import java.time.Instant;

@Service
public class UserService {

    private final PasswordEncoder passwordEncoder;
    private final UserRepository userRepository;
    private final AuthenticationManager authenticationManager;

    public UserService(PasswordEncoder passwordEncoder, UserRepository userRepository, AuthenticationManager authenticationManager) {
        this.passwordEncoder = passwordEncoder;
        this.userRepository = userRepository;
        this.authenticationManager = authenticationManager;
    }

    public User register(UserRegisterDto r) { //TODO user validation
        User user = new User();
        user.setFullName(r.getFullName());
        user.setUsername(r.getUsername());
        user.setPassword(passwordEncoder.encode(r.getPassword()));
        user.setCreatedAt(Instant.now());
        return userRepository.createUser(user); //TODO handle exception
    }

    public User authenticate(UserLoginDto user) {
        authenticationManager.authenticate(
                new UsernamePasswordAuthenticationToken(
                        user.getUsername(),
                        user.getPassword()
                )
        );
        return userRepository.getUserByUsername(user.getUsername());
    }
}
