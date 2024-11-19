package com.stringintech.security101.service;

import com.stringintech.security101.dto.UserLoginRequest;
import com.stringintech.security101.dto.UserRegisterRequest;
import com.stringintech.security101.exception.DuplicateUsernameException;
import com.stringintech.security101.exception.UserRegistrationException;
import com.stringintech.security101.model.User;
import com.stringintech.security101.repository.UserRepository;
import org.springframework.security.authentication.AuthenticationManager;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;

@Service
public class UserService {

    private final PasswordEncoder passwordEncoder;
    private final UserRepository userRepository;
    private final AuthenticationManager authenticationManager;
    private final TimeService timeService;

    public UserService(PasswordEncoder passwordEncoder, UserRepository userRepository,
                       AuthenticationManager authenticationManager, TimeService timeService) {
        this.passwordEncoder = passwordEncoder;
        this.userRepository = userRepository;
        this.authenticationManager = authenticationManager;
        this.timeService = timeService;
    }

    public User register(UserRegisterRequest r) {
        try {
            User user = new User();
            user.setFullName(r.getFullName());
            user.setUsername(r.getUsername());
            user.setPassword(passwordEncoder.encode(r.getPassword()));
            user.setCreatedAt(timeService.now());
            return userRepository.createUser(user);
        } catch (DuplicateUsernameException e) {
            throw e;
        } catch (Exception e) {
            throw new UserRegistrationException("Failed to register user", e);
        }
    }

    public User authenticate(UserLoginRequest user) {
        authenticationManager.authenticate(
                new UsernamePasswordAuthenticationToken(
                        user.getUsername(),
                        user.getPassword()
                )
        );
        return userRepository.getUserByUsername(user.getUsername());
    }
}
